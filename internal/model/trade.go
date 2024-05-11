package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Trade struct {
	ID        uint64
	TakerID   uint64
	MakerID   uint64
	Amount    decimal.Decimal
	Price     decimal.Decimal
	CreatedAt time.Time
}

type MatchResult struct {
	Taker        *Order
	FinishMakers []*Order
	Trades       []*Trade
}
