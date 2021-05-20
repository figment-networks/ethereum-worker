package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/figment-networks/ethereum-worker/api/conn"
	"github.com/figment-networks/ethereum-worker/structures"
	"github.com/figment-networks/indexing-engine/metrics"

	"go.uber.org/zap"
)

type Erc20API interface {
	TotalSupply(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (ts big.Int, err error)
	BalanceOf(ctx context.Context, bc *bind.BoundContract, tokenHolder common.Address, blockNumber uint64) (balance big.Int, err error)
	Name(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (name string, err error)
	Symbol(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (symbol string, err error)
	Decimals(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (res uint64, err error)
}

var (
	getAccountBalanceDuration *metrics.GroupObserver
)

// Client connecting to indexer-manager
type Client struct {
	serverApi Erc20API
	log       *zap.Logger
	ccm       *ContractCacheManager
	t         conn.EthereumTransport
	erc20ABI  abi.ABI
}

// NewClient is a indexer-manager Client constructor
func NewClient(log *zap.Logger, serverApi Erc20API, t conn.EthereumTransport, erc20ABI abi.ABI) *Client {
	return &Client{
		log:       log,
		t:         t,
		serverApi: serverApi,
		ccm:       NewContractCacheManager(),
		erc20ABI:  erc20ABI,
	}
}

func Init() {
	getAccountBalanceDuration = endpointDuration.WithLabels("getAccountBalance")
}

func (c *Client) LoadNetworkNames(ctx context.Context, name, address string) (err error) {
	cc := &ContractCache{BCC: c.t.GetBoundContractCaller(common.HexToAddress(address), c.erc20ABI)}
	if cc.Details, err = c.getERC20Details(ctx, cc.BCC.GetContract(), 0); err != nil {
		return fmt.Errorf("error calling getERC20Details: %w", err)
	}
	c.ccm.Set(address, name, cc)
	return nil
}

// GetAccountBalance returns account balance
func (c *Client) GetERC20AccountBalance(ctx context.Context, network, contract, address string, height uint64) ([]structures.Balance, error) {
	timer := metrics.NewTimer(getAccountBalanceDuration)
	defer timer.ObserveDuration()

	var (
		cc    *ContractCache
		found bool
	)

	if network != "" {
		cc, found = c.ccm.GetByNetwork(network)
	}
	if !found && address != "" {
		cc, found = c.ccm.GetByAddress(contract)
	}
	if !found {
		cc = &ContractCache{BCC: c.t.GetBoundContractCaller(common.HexToAddress(contract), c.erc20ABI)}
	}

	contractC := cc.BCC.GetContract()
	balance, err := c.serverApi.BalanceOf(ctx, contractC, common.HexToAddress(address), height)
	if err != nil {
		return nil, fmt.Errorf("error calling Balanceof: %w", err)
	}

	if !found {
		if cc.Details, err = c.getERC20Details(ctx, contractC, height); err != nil {
			return nil, fmt.Errorf("error calling getERC20Details: %w", err)
		}
		c.ccm.Set(contract, "", cc)
	}

	return []structures.Balance{{
		Values: structures.Values{
			Value: balance,
		},
		Details: cc.Details,
	}}, nil
}

func (c *Client) getERC20Details(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (det structures.Details, err error) {
	details := structures.Details{}
	if details.Name, err = c.serverApi.Name(ctx, bc, blockNumber); err != nil {
		return details, fmt.Errorf("error calling Name: %w", err)
	}

	if details.Symbol, err = c.serverApi.Symbol(ctx, bc, blockNumber); err != nil {
		return details, fmt.Errorf("error calling Symbol: %w", err)
	}

	if details.Decimals, err = c.serverApi.Decimals(ctx, bc, blockNumber); err != nil {
		return details, fmt.Errorf("error calling Decimals: %w", err)
	}
	return details, err
}
