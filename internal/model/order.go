package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID          uint64
	Side        OrderSide
	Type        OrderType
	Amount      decimal.Decimal // 所需數量
	Price       decimal.Decimal // 金額
	MatchAmount decimal.Decimal // 已成交數量
	CreatedAt   time.Time
}

// 剩餘所需數量
func (o *Order) GetLaveAmount() decimal.Decimal {
	return o.Amount.Sub(o.MatchAmount)
}
