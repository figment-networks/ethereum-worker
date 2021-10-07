package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/figment-networks/indexing-engine/metrics"
	"go.uber.org/zap"
)

var (
	getBalanceDuration *metrics.GroupObserver
	GetTotalSupplyDuration *metrics.GroupObserver
)

// GetBalance is http handler for GetBalance method
func (c *Connector) GetBalance(w http.ResponseWriter, req *http.Request) {
	timer := metrics.NewTimer(getBalanceDuration)
	defer timer.ObserveDuration()
	var (
		intHeight uint64
		err       error
	)
	enc := json.NewEncoder(w)
	height := req.URL.Query().Get("height")
	if height != "" {
		intHeight, err = strconv.ParseUint(height, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(ServiceError{Msg: "Invalid height param: " + err.Error()})
			return
		}
	}

	accountAddress := req.URL.Query().Get("accountAddress")
	if accountAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(ServiceError{Msg: "AccountAddress must be set"})
		return
	}

	network := req.URL.Query().Get("network")
	contractAddress := req.URL.Query().Get("contractAddress")
	if network == "" && contractAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(ServiceError{Msg: "Either network or contractAddress must be set"})
		return
	}

	ac, err := c.cli.GetERC20AccountBalance(req.Context(), network, contractAddress, accountAddress, intHeight)
	if err != nil {
		c.logger.Error("Error processing account request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(ServiceError{Msg: "Error processing account request"})
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = enc.Encode(ac); err != nil {
		c.logger.Error("Error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetTotalSupply is http handler for GetBalance method
func (c *Connector) GetTotalSupply(w http.ResponseWriter, req *http.Request) {
	timer := metrics.NewTimer(GetTotalSupplyDuration)
	defer timer.ObserveDuration()
	var (
		intHeight uint64
		err       error
	)
	enc := json.NewEncoder(w)
	height := req.URL.Query().Get("height")
	if height != "" {
		intHeight, err = strconv.ParseUint(height, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(ServiceError{Msg: "Invalid height param: " + err.Error()})
			return
		}
	}

	network := req.URL.Query().Get("network")
	contractAddress := req.URL.Query().Get("contractAddress")
	if network == "" && contractAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(ServiceError{Msg: "Either network or contractAddress must be set"})
		return
	}

	ac, err := c.cli.GetERC20TotalSupply(req.Context(), network, contractAddress, intHeight)
	if err != nil {
		c.logger.Error("Error processing account request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(ServiceError{Msg: "Error processing account request"})
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = enc.Encode(ac); err != nil {
		c.logger.Error("Error encoding response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
