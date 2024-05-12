package test

import (
	"Trading-Engine/internal/engine"
	"Trading-Engine/internal/model"
	"sync"
	"testing"
)

// 搓合效能測試 使用RBT
func BenchmarkMatchOrdersWithRBT(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := engine.NewOrderBook(engine.QueueTypePriceTree)
	b.Logf("orders len:%d", len(orders))
	for i := 0; i < b.N; i++ {
		r := orderBook.Match(orders[i])
		orderBook.UpdateFromMatchResult(*r)
	}
}

// 搓合效能測試 使用陣列
func BenchmarkMatchOrdersWithArray(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := engine.NewOrderBook(engine.QueueTypeArrayList)
	b.Logf("orders len:%d", len(orders))
	for i := 0; i < b.N; i++ {
		r := orderBook.Match(orders[i])
		orderBook.UpdateFromMatchResult(*r)
	}
}

func BenchmarkMatchOrdersWithRBTAndConcurrency(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := engine.NewOrderBook(engine.QueueTypePriceTree)
	b.Logf("orders len:%d", len(orders))
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(o model.Order) {
			defer wg.Done()
			orderBook.Mux.Lock()
			defer orderBook.Mux.Unlock()
			r := orderBook.Match(o)
			orderBook.UpdateFromMatchResult(*r)
		}(orders[i])
	}
	wg.Wait()
}

func BenchmarkMatchOrdersWithArrayAndConcurrency(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := engine.NewOrderBook(engine.QueueTypeArrayList)
	b.Logf("orders len:%d", len(orders))
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func(o model.Order) {
			defer wg.Done()
			orderBook.Mux.Lock()
			defer orderBook.Mux.Unlock()
			r := orderBook.Match(o)
			orderBook.UpdateFromMatchResult(*r)
		}(orders[i])
	}
	wg.Wait()
}
