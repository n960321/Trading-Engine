package model

import (
	"sync"
)

type OrderBook struct {
	buyQueue  OrderQueueior
	sellQueue OrderQueueior
	mu        sync.Mutex
}

type OrderQueueior interface {
	AddOrder(*Order)
	RemoveOrder(*Order)
	GetAllOrders() []*Order
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		buyQueue:  NewPriceTree(Buy),
		sellQueue: NewPriceTree(Sell),
	}
}

func (b *OrderBook) AddOrder(order *Order) {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch order.Side {
	case Buy:
		b.buyQueue.AddOrder(order)
	case Sell:
		b.sellQueue.AddOrder(order)
	}
}

func (b *OrderBook) GetAllBuyOrders() []*Order {
	return b.buyQueue.GetAllOrders()
}

func (b *OrderBook) GetAllSellOrders() []*Order {
	return b.sellQueue.GetAllOrders()
}

func (b *OrderBook) Match(order *Order) {
	switch order.Side {
	case Buy:
		// 找賣的本子有沒有
		// 有
		// 從最早的時間開始取order來配對
		// 沒有，加進買的本子等之後搓合
	case Sell:
		// 找買的本子有沒有
		// 有
		// 從最早的時間開始取order來配對
		// 沒有，加進賣的本子等之後搓合
	}
}

// func (b *OrderBook)
