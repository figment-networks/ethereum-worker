package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/figment-networks/ethereum-worker/api/conn/eth"
	"github.com/figment-networks/ethereum-worker/api/erc20"
	"github.com/figment-networks/ethereum-worker/client"
	"github.com/figment-networks/ethereum-worker/cmd/ethereum-worker-live/config"
	"github.com/figment-networks/ethereum-worker/cmd/ethereum-worker-live/logger"

	thttp "github.com/figment-networks/ethereum-worker/transport/http"

	"github.com/figment-networks/indexing-engine/health"
	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/metrics/prometheusmetrics"

	"go.uber.org/zap"
)

//go:embed abis/*
var abis embed.FS

type flags struct {
	configPath string
}

var configFlags = flags{}

func init() {
	flag.StringVar(&configFlags.configPath, "config", "", "path to config.json file")
	flag.Parse()
}

// Start runs ethereum-worker
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := getConfig(configFlags.configPath)
	if err != nil {
		log.Fatalf("error initializing config [ERR: %v]", err.Error())
	}

	if cfg.RollbarServerRoot == "" {
		cfg.RollbarServerRoot = "github.com/figment-networks/ethereum-worker"
	}
	rcfg := &logger.RollbarConfig{
		AppEnv:             cfg.AppEnv,
		RollbarAccessToken: cfg.RollbarAccessToken,
		RollbarServerRoot:  cfg.RollbarServerRoot,
		Version:            config.GitSHA,
	}

	if cfg.AppEnv == "development" || cfg.AppEnv == "local" {
		logger.Init("console", "debug", []string{"stderr"}, rcfg)
	} else {
		logger.Init("json", "info", []string{"stderr"}, rcfg)
	}

	logger.Info(config.IdentityString())
	defer logger.Sync()

	// Initialize metrics
	prom := prometheusmetrics.New()
	err = metrics.AddEngine(prom)
	if err != nil {
		logger.Error(err)
		return
	}
	err = metrics.Hotload(prom.Name())
	if err != nil {
		logger.Error(err)
	}

	tr := eth.NewEthTransport(cfg.EthereumAddress)
	if err := tr.Dial(ctx); err != nil {
		logger.Fatal("Error dialing ethereum", zap.String("ethereum_address", cfg.EthereumAddress), zap.Error(err))
		return
	}
	defer tr.Close(ctx)

	file, err := abis.ReadFile("abis/erc20abi.json")
	if err != nil {
		logger.Fatal("Error opening  erc20abi.json", zap.Error(err))
		return
	}
	erc20abi := &abi.ABI{}
	if err = json.Unmarshal(file, erc20abi); err != nil {
		logger.Fatal("Error opening  erc20abi.json", zap.Error(err))
		return
	}
	cl := client.NewClient(logger.GetLogger(), &erc20.ERC20Caller{}, tr, *erc20abi)
	client.Init()

	nNames := strings.Split(cfg.PredefinedNetworkNames, ";")
	for _, pair := range nNames {
		if !strings.ContainsAny(pair, ":") {
			logger.Fatal("PredefinedNetworkNames has to be in name:address;name:address;name:address format")
			return
		}
		network := strings.Split(pair, ":")
		logger.Info("Loading network: ", zap.String("name", network[0]), zap.String("address", network[1]), zap.String("node_address", cfg.EthereumAddress))
		if err = cl.LoadNetworkNames(ctx, network[0], network[1]); err != nil {
			logger.Fatal("Error loading network ", zap.Strings("config ", network), zap.Error(err))
		}
	}

	connector := thttp.NewConnector(cl, logger.GetLogger())
	mux := http.NewServeMux()

	connector.AttachToHandler(mux)
	mux.Handle("/metrics", metrics.Handler())

	monitor := &health.Monitor{}
	go monitor.RunChecks(ctx, cfg.HealthCheckInterval)
	monitor.AttachHttp(mux)

	handleHTTP(logger.GetLogger(), *cfg, mux)
}

func getConfig(path string) (cfg *config.Config, err error) {
	cfg = &config.Config{}
	if path != "" {
		if err := config.FromFile(path, cfg); err != nil {
			return nil, err
		}
	}

	if cfg.EthereumAddress != "" {
		return cfg, nil
	}

	if err := config.FromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func handleHTTP(l *zap.Logger, cfg config.Config, mux *http.ServeMux) {
	s := &http.Server{
		Addr:         cfg.Address + ":" + cfg.HTTPPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	l.Info("[HTTP] Listening on", zap.String("address", cfg.Address), zap.String("port", cfg.HTTPPort))
	if err := s.ListenAndServe(); err != nil {
		l.Error("[GRPC] Error while listening ", zap.String("address", cfg.Address), zap.String("port", cfg.HTTPPort), zap.Error(err))
	}
}
