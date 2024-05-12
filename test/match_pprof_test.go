package test

import (
	"Trading-Engine/internal/engine"
	"Trading-Engine/internal/model"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
)

// 用來看慢在哪
func Test_Matchpprof(t *testing.T) {
	orderNum := 50000
	orders := createOrder(int64(orderNum), model.OrderSideUnknow, model.OrderTypeUnknow, 1, 100, 1, 100)
	orderBook := engine.NewOrderBook(engine.QueueTypePriceTree)
	// orderBook := engine.NewOrderBook(engine.QueueTypeArrayList)
	f, _ := os.OpenFile("cpu-rbt.profile", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	var wg sync.WaitGroup
	wg.Add(orderNum)
	for i := 0; i < len(orders); i++ {
		go func(o model.Order) {
			defer wg.Done()
			orderBook.Mux.Lock()
			defer orderBook.Mux.Unlock()
			result := orderBook.Match(o)
			orderBook.UpdateFromMatchResult(*result)
		}(orders[i])
	}
	wg.Wait()
}
