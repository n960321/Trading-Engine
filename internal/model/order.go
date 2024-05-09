package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID          uint64          `json:"id" gorm:"column:id;primaryKey"`
	Side        OrderSide       `json:"side" gorm:"column:side;type:tinyint"`
	Type        OrderType       `json:"type" gorm:"column:type;type:tinyint"`
	Amount      decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(10,2)"`             // 所需數量
	Price       decimal.Decimal `json:"price" gorm:"column:price;type:decimal(10,2)"`               // 金額
	MatchAmount decimal.Decimal `json:"match_amount" gorm:"column:match_amount;type:decimal(10,2)"` // 已成交數量
	CreatedAt   time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   sql.NullTime    `json:"deleted_at" gorm:"column:deleted_at;"`
}

func (o *Order) TableName() string {
	return "orders"
}

// 剩餘所需數量
func (o *Order) GetLaveAmount() decimal.Decimal {
	return o.Amount.Sub(o.MatchAmount)
}
