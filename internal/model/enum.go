package model

// 訂單類型
type OrderType int

const (
	OrderTypeUnknow OrderType = iota // 未知
	OrderTypeLimit                   // 限價單
	OrderTypeMarket                  // 市價單
)

// 訂單走向
type OrderSide int

const (
	OrderSideUnknow OrderSide = iota // 未知
	OrderSideBuy                     // 買單
	OrderSideSell                    // 賣單
)
