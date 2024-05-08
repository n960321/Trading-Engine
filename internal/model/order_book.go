package model

import (
	"sync"

	"github.com/shopspring/decimal"
)

type OrderBook struct {
	Mux           sync.Mutex
	buyQueue      OrderQueueior
	sellQueue     OrderQueueior
	matchStrategy map[OrderType]MatchStrategy
}

func NewMatcher() map[OrderType]MatchStrategy {
	return map[OrderType]MatchStrategy{
		OrderTypeMarket: MarketStrategy{},
		OrderTypeLimit:  LimitStrategy{},
	}
}

type OrderQueueior interface {
	// 新增訂單
	AddOrder(*Order)
	// 移除訂單
	RemoveOrder(*Order)
	// 取得全部訂單
	GetAllOrders() []*Order
	// 從最優價格開始取得訂單到給入的price為止，且Total Amount等於或最後一筆訂單加上去會超過 Amount
	GetOrdersBetweenPirceWithAmount(price, laveAmount decimal.Decimal) []*Order
	// 從最優價格開始取訂單，直到Total Amount 等於或最後一筆訂單加上去會超過 Amount
	GetOrdersWithAmount(laveAmount decimal.Decimal) []*Order
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Mux:           sync.Mutex{},
		buyQueue:      NewPriceTree(OrderSideBuy),
		sellQueue:     NewPriceTree(OrderSideSell),
		matchStrategy: NewMatcher(),
	}
}

// 將訂單加入訂單簿但不進行搓合
func (b *OrderBook) AddOrder(order *Order) {
	switch order.Side {
	case OrderSideBuy:
		b.buyQueue.AddOrder(order)
	case OrderSideSell:
		b.sellQueue.AddOrder(order)
	}
}

func (b *OrderBook) CancelOrder(order *Order) {
	switch order.Side {
	case OrderSideBuy:
		b.buyQueue.RemoveOrder(order)
	case OrderSideSell:
		b.sellQueue.RemoveOrder(order)
	}

}

// 取得所有掛買單
func (b *OrderBook) GetAllBuyOrders() []*Order {
	return b.buyQueue.GetAllOrders()
}

// 取得所有掛賣單
func (b *OrderBook) GetAllSellOrders() []*Order {
	return b.sellQueue.GetAllOrders()
}

// 將訂單與訂單簿進行撮合
func (b *OrderBook) Match(order *Order) *MatchResult {
	b.Mux.Lock()
	defer b.Mux.Unlock()
	var matchResult *MatchResult
	if strategy, exist := b.matchStrategy[order.Type]; exist {
		matchResult = strategy.MatchOrder(b.buyQueue, b.sellQueue, order)
	}

	// 從訂單簿中移除已經搓合完成的makers
	for _, maker := range matchResult.FinishMakers {
		switch maker.Side {
		case OrderSideBuy:
			b.buyQueue.RemoveOrder(maker)
		case OrderSideSell:
			b.sellQueue.RemoveOrder(maker)
		}
	}

	return matchResult
}