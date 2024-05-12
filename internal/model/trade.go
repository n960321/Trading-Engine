package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Trade struct {
	gorm.Model
	TakerID uint            `json:"TakerID" gorm:"column:taker_id"`
	MakerID uint            `json:"MakerID" gorm:"column:maker_id"`
	Amount  decimal.Decimal `json:"Amount" gorm:"column:amount;type:decimal(10,2)"`
	Price   decimal.Decimal `json:"Price" gorm:"column:price;type:decimal(10,2)"`
}

type MatchResult struct {
	Taker           Order
	FinishMakers    []Order
	IncompleteMaker Order
	Trades          []Trade
}
