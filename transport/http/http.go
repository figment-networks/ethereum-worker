package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/figment-networks/ethereum-worker/structures"
	"go.uber.org/zap"
)

type retrieveClienter interface {
	GetERC20AccountBalance(ctx context.Context, network, contract, address string, height uint64) ([]structures.Balance, error)
	GetERC20TotalSupply(ctx context.Context, network, contract string, height uint64) ([]structures.Balance, error)
}

// Connector is main HTTP connector for manager
type Connector struct {
	cli    retrieveClienter
	logger *zap.Logger
}

// NewConnector is  Connector constructor
func NewConnector(cli retrieveClienter, logger *zap.Logger) *Connector {
	getBalanceDuration = endpointDuration.WithLabels("getBalance")
	GetTotalSupplyDuration = endpointDuration.WithLabels("getTotalSupply")
	return &Connector{cli, logger}
}

// AttachToHandler attaches handlers to http server's mux
func (c *Connector) AttachToHandler(mux *http.ServeMux) {
	mux.HandleFunc("/getBalance", c.GetBalance)
	mux.HandleFunc("/getTotalSupply", c.GetTotalSupply)
}

// ServiceError structure as formated error
type ServiceError struct {
	Status int         `json:"status"`
	Msg    interface{} `json:"error"`
}

func (ve ServiceError) Error() string {
	return fmt.Sprintf("Bad Request: %s", ve.Msg)
}
