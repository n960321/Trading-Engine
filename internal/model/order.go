package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Side        OrderSide       `json:"Side" gorm:"column:side;type:tinyint"`
	Type        OrderType       `json:"Type" gorm:"column:type;type:tinyint"`
	Amount      decimal.Decimal `json:"Amount" gorm:"column:amount;type:decimal(10,2)"`            // 所需數量
	Price       decimal.Decimal `json:"Price" gorm:"column:price;type:decimal(10,2)"`              // 金額
	MatchAmount decimal.Decimal `json:"MatchAmount" gorm:"column:match_amount;type:decimal(10,2)"` // 已成交數量
	Compeleted  bool            `json:"Completed" gorm:"column:completed;type:tinyint(1)"`         // 是否已完成
}

func (o *Order) TableName() string {
	return "orders"
}

// 剩餘所需數量
func (o *Order) GetLaveAmount() decimal.Decimal {
	return o.Amount.Sub(o.MatchAmount)
}
