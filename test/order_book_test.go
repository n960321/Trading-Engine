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

	sellOrders := createOrder(RandomInt(10000, 10000000), model.OrderSideSell, 1, 100, 10, 200)
	buyOrders := createOrder(RandomInt(10000, 10000000), model.OrderSideBuy, 1, 100, 10, 200)

	sellResult := make([]*model.Order, len(sellOrders))
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
	buyResult := make([]*model.Order, len(buyOrders))
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
	type args struct {
		orders []*model.Order
		result []*model.Order
		side   model.OrderSide
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "加入多筆掛賣訂單 確認order book 中掛賣的順序是對的",
			args: args{
				side:   model.OrderSideSell,
				orders: sellOrders,
				result: sellResult,
			},
		},
		{
			name: "加入多筆掛買訂單 確認order book 中掛買的順序是對的",
			args: args{
				side:   model.OrderSideBuy,
				orders: buyOrders,
				result: buyResult,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderBook := model.NewOrderBook()
			for _, order := range tt.args.orders {
				orderBook.AddOrder(order)
			}

			var orderInOrderBook []*model.Order

			switch tt.args.side {
			case model.OrderSideBuy:
				orderInOrderBook = orderBook.GetAllBuyOrders()
			case model.OrderSideSell:
				orderInOrderBook = orderBook.GetAllSellOrders()
			}
			t.Logf("Total orders : %d", len(orderInOrderBook))
			if len(orderInOrderBook) != len(tt.args.result) {
				t.Errorf("AddOrders(%s) result lengh not equal, len orderInOrderBook: %d , len result: %d", tt.name, len(orderInOrderBook), len(tt.args.result))
			}

			for i, order := range orderInOrderBook {
				if order.ID != tt.args.result[i].ID {
					t.Errorf("AddOrder(%s) the sort is wrong!", tt.name)
					t.Failed()
				}
			}
		})
	}
}


func Test_MatchOrders(t *testing.T) {
	orders := createOrder(RandomInt(1000000, 1000000), model.OrderSideUnknow, 1, 100, 1, 100)
	type args struct {
		orders []*model.Order
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "加入一百萬筆訂單，嘗試撮合看可以幾秒內完成",
			args: args{
				orders: orders,
			},
		},
	}

	for _, tt := range tests {
		t.Logf("orders len: %d",len(tt.args.orders))
		t.Run(tt.name, func(t *testing.T) {
			orderBook := model.NewOrderBook()
			for _, order := range orders {
				orderBook.Match(order)
			}
		})
	}
}

// TODO : 需要補一個測試是驗證撮合的正確性！


func createOrder(q int64, side model.OrderSide, priceMin, priceMax, amountMin, amountMax int64) []*model.Order {
	orders := make([]*model.Order, 0, q)
	for i := int64(0); i < q; i++ {
		curSide := model.OrderSideUnknow
		if side == model.OrderSideUnknow {
			curSide = model.OrderSide(RandomInt(1, 2))
		}
		orders = append(orders, &model.Order{
			ID:        uint64(i),
			Side:      curSide,
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
