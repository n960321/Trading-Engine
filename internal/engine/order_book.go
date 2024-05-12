package engine

import (
	"Trading-Engine/internal/model"
	"sync"

	"github.com/shopspring/decimal"
)

type OrderBook struct {
	Mux           sync.Mutex
	buyQueue      OrderQueueior
	sellQueue     OrderQueueior
	matchStrategy map[model.OrderType]MatchStrategy
}

func NewMatcher() map[model.OrderType]MatchStrategy {
	return map[model.OrderType]MatchStrategy{
		model.OrderTypeMarket: MarketStrategy{},
		model.OrderTypeLimit:  LimitStrategy{},
	}
}

type OrderQueueior interface {
	// 新增訂單 如果已經存在要覆蓋
	AddOrder(model.Order)
	// 移除訂單
	RemoveOrder(model.Order)
	// 取得全部訂單
	GetAllOrders() []model.Order
	// 從最優價格開始取得訂單到給入的price為止，且Total Amount等於或最後一筆訂單加上去會超過 Amount
	GetOrdersBetweenPirceWithAmount(price, laveAmount decimal.Decimal) []model.Order
	// 從最優價格開始取訂單，直到Total Amount 等於或最後一筆訂單加上去會超過 Amount
	GetOrdersWithAmount(laveAmount decimal.Decimal) []model.Order
}

var queueTypeMap = map[model.QueueType](func(side model.OrderSide) OrderQueueior){
	QueueTypeArrayList: NewQueueByList,
	QueueTypePriceTree: NewRBTPriceTree,
}

func NewOrderBook(queueType model.QueueType) *OrderBook {
	return &OrderBook{
		Mux:           sync.Mutex{},
		buyQueue:      queueTypeMap[queueType](model.OrderSideBuy),
		sellQueue:     queueTypeMap[queueType](model.OrderSideSell),
		matchStrategy: NewMatcher(),
	}
}

func (b *OrderBook) AddOrder(order model.Order) {
	switch order.Side {
	case model.OrderSideBuy:
		b.buyQueue.AddOrder(order)
	case model.OrderSideSell:
		b.sellQueue.AddOrder(order)
	}
}

func (b *OrderBook) RemoveOrder(order model.Order) {
	switch order.Side {
	case model.OrderSideBuy:
		b.buyQueue.RemoveOrder(order)
	case model.OrderSideSell:
		b.sellQueue.RemoveOrder(order)
	}
}

// 取得所有掛買單
func (b *OrderBook) GetAllBuyOrders() []model.Order {
	return b.buyQueue.GetAllOrders()
}

// 取得所有掛賣單
func (b *OrderBook) GetAllSellOrders() []model.Order {
	return b.sellQueue.GetAllOrders()
}

// 將訂單與訂單簿進行撮合
func (b *OrderBook) Match(order model.Order) *model.MatchResult {
	var matchResult model.MatchResult
	if strategy, exist := b.matchStrategy[order.Type]; exist {
		matchResult = strategy.MatchOrder(b.buyQueue, b.sellQueue, order)
	}

	return &matchResult
}

func (b *OrderBook) UpdateFromMatchResult(m model.MatchResult) {
	// 如果完成 taker 要移除 ,沒的話要更新 -> 其實用同一個就好
	if m.Taker.Type == model.OrderTypeLimit {
		if m.Taker.Compeleted {
			b.RemoveOrder(m.Taker)
		} else {
			b.AddOrder(m.Taker)
		}
	}
	// imcomplete 要更新
	if m.IncompleteMaker.ID > 0 {
		b.AddOrder(m.IncompleteMaker)
	}

	// finish maker 要移除
	for _, maker := range m.FinishMakers {
		b.RemoveOrder(maker)
	}

}
