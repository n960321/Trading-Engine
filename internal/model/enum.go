package model

// 訂單操作
type OrderAction int

const (
	OrderActionUnknow OrderAction = iota // 未知
	OrderActionCreate                    // 創建掛單
	OrderActionCancel                    // 取消掛單
)

// 訂單走向
type OrderSide int

const (
	OrderSideUnknow OrderSide = iota // 未知
	OrderSideBuy                     // 買單
	OrderSideSell                    // 賣單
)
