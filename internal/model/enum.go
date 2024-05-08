package model

// 訂單走向
type OrderSide int

const (
	OrderSideUnknow OrderSide = iota // 未知
	OrderSideBuy                     // 買單
	OrderSideSell                    // 賣單
)
