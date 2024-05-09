package test

import (
	"Trading-Engine/internal/model"
	"os"
	"runtime/pprof"
	"testing"
)

// 用來看慢在哪
func Test_Matchpprof(t *testing.T) {
	orders := createOrder(int64(1000000), model.OrderSideUnknow, model.OrderTypeLimit, 1, 100, 1, 100)
	f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	orderBook := model.NewOrderBook()
	for i := 0; i < len(orders); i++ {
		orderBook.Match(orders[i])
	}
}
