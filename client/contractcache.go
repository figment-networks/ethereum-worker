package client

import (
	"strings"
	"sync"

	"github.com/figment-networks/ethereum-worker/api/conn"
	"github.com/figment-networks/ethereum-worker/structures"
)

type ContractCache struct {
	BCC     conn.BoundContractCaller
	Details structures.Details
}

type ContractCacheManager struct {
	l          sync.RWMutex
	addressMap map[string]*ContractCache
	networkMap map[string]*ContractCache
}

func NewContractCacheManager() *ContractCacheManager {
	return &ContractCacheManager{addressMap: make(map[string]*ContractCache), networkMap: make(map[string]*ContractCache)}
}

func (cc *ContractCacheManager) GetByAddress(address string) (*ContractCache, bool) {
	cc.l.RLock()
	defer cc.l.RUnlock()
	b, ok := cc.addressMap[address]
	return b, ok
}

func (cc *ContractCacheManager) GetByNetwork(network string) (*ContractCache, bool) {
	cc.l.RLock()
	defer cc.l.RUnlock()
	b, ok := cc.networkMap[strings.ToLower(network)]
	return b, ok
}

func (cc *ContractCacheManager) Set(address, network string, contract *ContractCache) {
	cc.l.Lock()
	defer cc.l.Unlock()

	cc.addressMap[address] = contract
	if network != "" {
		cc.networkMap[strings.ToLower(network)] = contract
	}
}
