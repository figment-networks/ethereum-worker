package structures

import "math/big"

type Balance struct {
	Values  Values  `json:"values"`
	Details Details `json:"details"`
}

type Details struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals uint64 `json:"decimals"`
}

type Values struct {
	Value big.Int `json:"value"`
	Type  string  `json:"type"`
}
