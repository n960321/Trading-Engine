package model

// 訂單操作
type OrderAction int

const (
	Cancel OrderAction = iota // 取消掛單
	Create                    // 創建掛單
)

// 訂單走向
type OrderSide int

const (
	Buy  OrderSide = iota // 買單
	Sell                  // 賣單
)
