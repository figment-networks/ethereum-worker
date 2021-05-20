package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/figment-networks/ethereum-worker/api/conn"
)

type EthTransport struct {
	C   *ethclient.Client
	Url string
}

func NewEthTransport(url string) *EthTransport {
	return &EthTransport{Url: url}
}

func (et *EthTransport) Dial(ctx context.Context) (err error) {
	et.C, err = ethclient.DialContext(ctx, et.Url)
	return err
}

func (et *EthTransport) Close(ctx context.Context) {
	et.C.Close()
}

func (et *EthTransport) GetBoundContractCaller(address common.Address, a abi.ABI) conn.BoundContractCaller {
	return &BoundContractC{address: address, abi: a, ET: et}
}

type BoundContractC struct {
	address common.Address
	abi     abi.ABI
	ET      *EthTransport
}

func (bcc *BoundContractC) GetContract() *bind.BoundContract {
	return bind.NewBoundContract(bcc.address, bcc.abi, bcc.ET.C, nil, nil) // (lukanus): it's just a structure

}
