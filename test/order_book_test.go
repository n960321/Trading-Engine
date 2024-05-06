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
	
	sellOrders := createOrder(RandomInt(10000, 10000000), model.Sell, 1, 100, 10, 200)
	buyOrders := createOrder(RandomInt(10000, 10000000), model.Buy, 1, 100, 10, 200)

	sellResult := make([]*model.Order, len(sellOrders))
	copy(sellResult, sellOrders)

	sort.Slice(sellResult, func(i, j int) bool {
		r := sellResult[i].Price.Compare(sellResult[j].Price)
		if r == 0 {
			return !sellResult[i].CreatedAt.After(sellResult[j].CreatedAt)
		} else if r == -1 {
			return false
		}
		return true
	})
	buyResult := make([]*model.Order, len(buyOrders))
	copy(buyResult, buyOrders)

	sort.Slice(buyResult, func(i, j int) bool {
		r := buyResult[i].Price.Compare(buyResult[j].Price)
		if r == 0 {
			return !buyResult[i].CreatedAt.After(buyResult[j].CreatedAt)
		} else if r == -1 {
			return true
		}
		return false
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
			name: "加入多筆掛買訂單 確認order book 中掛買的順序是對的",
			args: args{
				side:   model.Buy,
				orders: buyOrders,
				result: buyResult,
			},
		},
		{
			name: "加入多筆掛買訂單 確認order book 中掛賣的順序是對的",
			args: args{
				side:   model.Sell,
				orders: sellOrders,
				result: sellResult,
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
			case model.Buy:
				orderInOrderBook = orderBook.GetAllBuyOrders()
			case model.Sell:
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

func createOrder(q int64, side model.OrderSide, priceMin, priceMax, amountMin, amountMax int64) []*model.Order {
	orders := make([]*model.Order, 0, q)
	for i := int64(0); i < q; i++ {
		orders = append(orders, &model.Order{
			ID:        uint64(i),
			Side:      side,
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
