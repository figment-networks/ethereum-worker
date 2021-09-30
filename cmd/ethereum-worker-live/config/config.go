package config

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	Name      = "ethereum-worker"
	Version   string
	GitSHA    string
	Timestamp string
)

const (
	modeDevelopment = "development"
	modeProduction  = "production"
)

// Config holds the configuration data
type Config struct {
	AppEnv string `json:"app_env" envconfig:"APP_ENV" default:"development"`

	Address  string `json:"address" envconfig:"ADDRESS" default:"0.0.0.0"`
	HTTPPort string `json:"http_port" envconfig:"HTTP_PORT" default:"8097"`

	EthereumAddress        string `json:"ethereum_address" envconfig:"ETHEREUM_ADDRESS" default:"http://host.docker.internal:8545"` // <---- probably needs to be fixed
	PredefinedNetworkNames string `json:"predefined_network_named" envconfig:"PREDEFINED_NETWORK_NAMES" default:"skale:0x00c83aeCC790e8a4453e5dD3B0B4b3680501a7A7"`

	// Rollbar
	RollbarAccessToken string `json:"rollbar_access_token" envconfig:"ROLLBAR_ACCESS_TOKEN"`
	RollbarServerRoot  string `json:"rollbar_server_root" envconfig:"ROLLBAR_SERVER_ROOT" default:"github.com/figment-networks/account-service"`

	HealthCheckInterval time.Duration `json:"health_check_interval" envconfig:"HEALTH_CHECK_INTERVAL" default:"10s"`
}

// FromFile reads the config from a file
func FromFile(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// FromEnv reads the config from environment variables
func FromEnv(config *Config) error {
	return envconfig.Process("", config)
}
