package test

import (
	"Trading-Engine/internal/model"
	"sync"
	"testing"
)

// 搓合效能測試 使用RBT
func BenchmarkMatchOrdersWithRBT(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := model.NewOrderBook(model.QueueTypePriceTree)
	b.Logf("orders len:%d", len(orders))
	for i := 0; i < b.N; i++ {
		orderBook.Match(&orders[i])
	}
}

// 搓合效能測試 使用陣列
func BenchmarkMatchOrdersWithArray(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := model.NewOrderBook(model.QueueTypeArrayList)
	b.Logf("orders len:%d", len(orders))
	for i := 0; i < b.N; i++ {
		orderBook.Match(&orders[i])
	}
}

func BenchmarkMatchOrdersWithArrayAndConcurrency(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := model.NewOrderBook(model.QueueTypeArrayList)
	b.Logf("orders len:%d", len(orders))
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			orderBook.Match(&orders[i])
		}()
	}
	wg.Wait()
}

func BenchmarkMatchOrdersWithRBTAndConcurrency(b *testing.B) {
	orders := createOrder(int64(b.N), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := model.NewOrderBook(model.QueueTypePriceTree)
	b.Logf("orders len:%d", len(orders))
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			orderBook.Match(&orders[i])
		}()
	}
	wg.Wait()
}