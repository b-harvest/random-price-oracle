package main

import (
	"time"
)

const (
	CoinSymbolKey = "symbol"
	CoinPriceKey  = "price"
)

type Coin struct {
	Symbol    string    `bson:"symbol"`
	Price     float64   `bson:"price"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type GetPricesRequest struct {
	Symbols Strings `query:"symbols"`
}

type GetPricesResponse struct {
	Coins map[string]GetPricesResponseCoin `json:"coins"`
}

type GetPricesResponseCoin struct {
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"updatedAt"`
}
