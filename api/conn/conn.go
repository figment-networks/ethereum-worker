package conn

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var ErrEmptyResponse = errors.New("Returned Empty Response (Reverted 0x)")

type BoundContractCaller interface {
	GetContract() *bind.BoundContract
}

type EthereumTransport interface {
	Dial(ctx context.Context) (err error)
	Close(ctx context.Context)
	GetBoundContractCaller(address common.Address, a abi.ABI) BoundContractCaller
}
