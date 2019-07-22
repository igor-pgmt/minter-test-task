package server

import (
	"math/big"

	"github.com/igor-pgmt/minter-test-task/client"
)

type Response struct {
	Result struct {
		Sum          *big.Int             `json:"sum,omitempty"`
		Transactions []client.Transaction `json:"transactions,omitempty"`
	} `json:"result"`
}

type ResponseTime struct {
	Result struct {
		Transactions []Transactions `json:"transactions,omitempty"`
	} `json:"result"`
}

type Transactions struct {
	From         string               `json:"from,omitempty"`
	Sum          *big.Int             `json:"sum,omitempty"`
	Transactions []client.Transaction `json:"transactions,omitempty"`
}
