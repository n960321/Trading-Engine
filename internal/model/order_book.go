package model

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type OrderBook struct {
	Mux       sync.Mutex
	buyQueue  OrderQueueior
	sellQueue OrderQueueior
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
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Mux:       sync.Mutex{},
		buyQueue:  NewPriceTree(OrderSideBuy),
		sellQueue: NewPriceTree(OrderSideSell),
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
	var makers []*Order
	switch order.Side {
	case OrderSideBuy:
		makers = b.sellQueue.GetOrdersBetweenPirceWithAmount(order.Price, order.GetLaveAmount())
	case OrderSideSell:
		makers = b.buyQueue.GetOrdersBetweenPirceWithAmount(order.Price, order.GetLaveAmount())
	}

	if len(makers) == 0 {
		// log something
		b.AddOrder(order)
		return &MatchResult{Taker: order}
	}
	trades := make([]*Trade, 0)
	for _, maker := range makers {
		// 當 taker 剩餘Amount 大於等於 maker的LaveAmount
		if order.GetLaveAmount().Compare(maker.GetLaveAmount()) >= 0 {
			trades = append(trades, &Trade{
				TakerID:   order.ID,
				MakerID:   maker.ID,
				TakerSide: order.Side,
				Amount:    maker.GetLaveAmount(),
				Price:     maker.Price,
				CreatedAt: time.Now(),
			})

			order.MatchAmount.Add(maker.GetLaveAmount())
			maker.MatchAmount.Add(maker.GetLaveAmount())

		} else {
			// 當 taker 剩餘Amount 小於 maker的LaveAmount
			trades = append(trades, &Trade{
				TakerID:   order.ID,
				MakerID:   maker.ID,
				TakerSide: order.Side,
				Amount:    order.GetLaveAmount(),
				Price:     maker.Price,
				CreatedAt: time.Now(),
			})

			order.MatchAmount.Add(order.GetLaveAmount())
			maker.MatchAmount.Add(order.GetLaveAmount())
			// 由於最後一個還有為搓合的數量，所以移除掉
			makers = makers[:len(makers)-1]
		}
	}
	// taker 剩餘Amount 不是為零 則需要再放進訂單簿中等待搓合
	if !order.GetLaveAmount().Equal(decimal.Zero) {
		b.AddOrder(order)
	}

	// 從訂單簿中移除已經搓合完成的makers
	for _, maker := range makers {
		switch maker.Side {
		case OrderSideBuy:
			b.buyQueue.RemoveOrder(maker)
		case OrderSideSell:
			b.sellQueue.RemoveOrder(maker)
		}
	}

	return &MatchResult{
		Taker:        order,
		FinishMakers: makers,
		Trades:       trades,
	}

}
