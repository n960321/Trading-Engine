package test

import (
	"Trading-Engine/internal/model"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// 加入多筆訂單 確認order book 中的順序是對的
func Test_AddOrders(t *testing.T) {

	getTestCase := func(min, max int64) (sellOrders, buyOrders, sellResult, buyResult []model.Order) {
		sellOrders = createOrder(RandomInt(min, max), model.OrderSideSell, model.OrderTypeLimit, 1, 100, 10, 200)
		buyOrders = createOrder(RandomInt(min, max), model.OrderSideBuy, model.OrderTypeLimit, 1, 100, 10, 200)

		sellResult = make([]model.Order, len(sellOrders))
		copy(sellResult, sellOrders)

		sort.Slice(sellResult, func(i, j int) bool {
			r := sellResult[i].Price.Compare(sellResult[j].Price)
			if r == 0 {
				return !sellResult[i].CreatedAt.After(sellResult[j].CreatedAt)
			} else if r == -1 {
				return true
			}
			return false
		})
		buyResult = make([]model.Order, len(buyOrders))
		copy(buyResult, buyOrders)

		sort.Slice(buyResult, func(i, j int) bool {
			r := buyResult[i].Price.Compare(buyResult[j].Price)
			if r == 0 {
				return !buyResult[i].CreatedAt.After(buyResult[j].CreatedAt)
			} else if r == -1 {
				return false
			}
			return true
		})
		return sellOrders, buyOrders, sellResult, buyResult
	}

	type args struct {
		getTestCaseFunc func(min, max int64) (sellOrders, buyOrders, sellResult, buyResult []model.Order)
		side            model.OrderSide
		queueType       model.QueueType
		min             int64
		max             int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "RBT - 加入多筆掛賣 and 掛買 確認order book 中掛賣的順序是對的",
			args: args{
				side:            model.OrderSideSell,
				getTestCaseFunc: getTestCase,
				queueType:       model.QueueTypePriceTree,
				min:             10000,
				max:             1000000,
			},
		},
		{
			name: "ArrayList - 加入多筆掛買訂單 確認order book 中掛買的順序是對的",
			args: args{
				side:            model.OrderSideBuy,
				getTestCaseFunc: getTestCase,
				queueType:       model.QueueTypeArrayList,
				min:             1000,
				max:             100000,
			},
		},
	}

	for _, tt := range tests {
		sellOrders, buyOrders, sellResult, buyResult := tt.args.getTestCaseFunc(tt.args.min, tt.args.max)
		t.Run(tt.name+"with 賣", func(t *testing.T) {
			orderBook := model.NewOrderBook(tt.args.queueType)
			for _, order := range sellOrders {
				orderBook.AddOrder(&order)
			}
			orderInOrderBook := orderBook.GetAllSellOrders()

			t.Logf("Total orders : %d", len(orderInOrderBook))
			if len(orderInOrderBook) != len(sellResult) {
				t.Errorf("AddOrders(%s) result lengh not equal, len orderInOrderBook: %d , len result: %d", tt.name, len(orderInOrderBook), len(sellResult))
			}

			for i, order := range orderInOrderBook {
				if order.ID != sellResult[i].ID {
					t.Errorf("AddOrder(%s) the sort is wrong!", tt.name)
					t.Failed()
				}
			}
		})

		t.Run(tt.name+"with 買", func(t *testing.T) {
			orderBook := model.NewOrderBook(tt.args.queueType)
			for _, order := range buyOrders {
				orderBook.AddOrder(&order)
			}
			orderInOrderBook := orderBook.GetAllBuyOrders()

			t.Logf("Total orders : %d", len(orderInOrderBook))
			if len(orderInOrderBook) != len(buyResult) {
				t.Errorf("AddOrders(%s) result lengh not equal, len orderInOrderBook: %d , len result: %d", tt.name, len(orderInOrderBook), len(buyResult))
			}

			for i, order := range orderInOrderBook {
				if order.ID != buyResult[i].ID {
					t.Errorf("AddOrder(%s) the sort is wrong!", tt.name)
					t.Failed()
				}
			}
		})
	}
}

// 加入多筆訂單後，取消某筆訂單，確認是否還在OrderBook中
func Test_CancelOrder(t *testing.T) {
	orders := createOrder(RandomInt(100, 1000), model.OrderSideUnknow, model.OrderTypeLimit, 1, 100, 1, 100)
	cancelOrder := orders[RandomInt(0, int64(len(orders)))]

	type args struct {
		orders      []model.Order
		cancelOrder model.Order
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "加入多筆訂單，隨機一筆然後刪除",
			args: args{
				orders:      orders,
				cancelOrder: cancelOrder,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderBook := model.NewOrderBook(model.QueueTypePriceTree)
			for _, order := range tt.args.orders {
				orderBook.AddOrder(&order)
			}

			orderBook.CancelOrder(&tt.args.cancelOrder)

			curOrders := orderBook.GetAllBuyOrders()
			curOrders = append(curOrders, orderBook.GetAllSellOrders()...)

			for _, order := range curOrders {
				if order.ID == tt.args.cancelOrder.ID {
					t.Errorf("The order has not been canceled.")
					break
				}
			}
		})
	}

}

// TODO : 需要補一個測試是驗證撮合的正確性 -> 只驗搓合的正確性，訂單簿的順序不在此測試範圍
// 跑得久沒關係 可以用最暴力的方法且一定正確的方式，當對照組，然後要測的當實驗組

// func

func createOrder(q int64, side model.OrderSide, orderType model.OrderType, priceMin, priceMax, amountMin, amountMax int64) []model.Order {
	orders := make([]model.Order, 0, q)
	for i := int64(0); i < q; i++ {
		curSide := side
		if curSide == model.OrderSideUnknow {
			curSide = model.OrderSide(RandomInt(1, 2))
		}
		curType := orderType

		if curType == model.OrderTypeUnknow {
			curType = model.OrderType(RandomInt(1, 2))
		}

		orders = append(orders, model.Order{
			ID:        uint64(i),
			Side:      curSide,
			Type:      curType,
			Amount:    decimal.NewFromInt(RandomInt(amountMin, amountMax)),
			Price:     decimal.NewFromInt(RandomInt(priceMin, priceMax)),
			CreatedAt: time.Now(),
		})
	}
	return orders

}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
