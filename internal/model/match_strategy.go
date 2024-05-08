package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type MatchStrategy interface {
	MatchOrder(buyQueue, sellQueue OrderQueueior, taker *Order) *MatchResult
}

// 限價單撮合
type LimitStrategy struct{}

func (l LimitStrategy) MatchOrder(buyQueue, sellQueue OrderQueueior, taker *Order) *MatchResult {
	var makers []*Order
	switch taker.Side {
	case OrderSideBuy:
		makers = sellQueue.GetOrdersBetweenPirceWithAmount(taker.Price, taker.GetLaveAmount())
	case OrderSideSell:
		makers = buyQueue.GetOrdersBetweenPirceWithAmount(taker.Price, taker.GetLaveAmount())
	}

	if len(makers) == 0 {
		if taker.Side == OrderSideBuy {
			buyQueue.AddOrder(taker)
		} else {
			sellQueue.AddOrder(taker)
		}
		return &MatchResult{Taker: taker}
	}
	trades := make([]*Trade, 0)
	for _, maker := range makers {
		// 當 taker 剩餘Amount 大於等於 maker的LaveAmount
		if taker.GetLaveAmount().Compare(maker.GetLaveAmount()) >= 0 {
			trades = append(trades, &Trade{
				TakerID:   taker.ID,
				MakerID:   maker.ID,
				TakerSide: taker.Side,
				Amount:    maker.GetLaveAmount(),
				Price:     maker.Price,
				CreatedAt: time.Now(),
			})

			taker.MatchAmount.Add(maker.GetLaveAmount())
			maker.MatchAmount.Add(maker.GetLaveAmount())

		} else {
			// 當 taker 剩餘Amount 小於 maker的LaveAmount
			trades = append(trades, &Trade{
				TakerID:   taker.ID,
				MakerID:   maker.ID,
				TakerSide: taker.Side,
				Amount:    taker.GetLaveAmount(),
				Price:     maker.Price,
				CreatedAt: time.Now(),
			})

			taker.MatchAmount.Add(taker.GetLaveAmount())
			maker.MatchAmount.Add(taker.GetLaveAmount())
			// 由於最後一個還有為搓合的數量，所以移除掉
			makers = makers[:len(makers)-1]
		}
	}
	// taker 剩餘Amount 不是為零 則需要再放進訂單簿中等待搓合
	if !taker.GetLaveAmount().Equal(decimal.Zero) {
		if taker.Side == OrderSideBuy {
			buyQueue.AddOrder(taker)
		} else {
			sellQueue.AddOrder(taker)
		}
	}

	return &MatchResult{
		Taker:        taker,
		FinishMakers: makers,
		Trades:       trades,
	}
}

// 市價單撮合
type MarketStrategy struct{}

func (l MarketStrategy) MatchOrder(buyQueue, sellQueue OrderQueueior, taker *Order) *MatchResult {
	var makers []*Order
	switch taker.Side {
	case OrderSideBuy:
		makers = sellQueue.GetOrdersWithAmount(taker.GetLaveAmount())
	case OrderSideSell:
		makers = buyQueue.GetOrdersWithAmount(taker.GetLaveAmount())
	}

	if len(makers) == 0 {
		return &MatchResult{Taker: taker}
	}
	trades := make([]*Trade, 0)
	for _, maker := range makers {
		// 當 taker 剩餘Amount 大於等於 maker的LaveAmount
		if taker.GetLaveAmount().Compare(maker.GetLaveAmount()) >= 0 {
			trades = append(trades, &Trade{
				TakerID:   taker.ID,
				MakerID:   maker.ID,
				TakerSide: taker.Side,
				Amount:    maker.GetLaveAmount(),
				Price:     maker.Price,
				CreatedAt: time.Now(),
			})

			taker.MatchAmount.Add(maker.GetLaveAmount())
			maker.MatchAmount.Add(maker.GetLaveAmount())

		} else {
			// 當 taker 剩餘Amount 小於 maker的LaveAmount
			trades = append(trades, &Trade{
				TakerID:   taker.ID,
				MakerID:   maker.ID,
				TakerSide: taker.Side,
				Amount:    taker.GetLaveAmount(),
				Price:     maker.Price,
				CreatedAt: time.Now(),
			})

			taker.MatchAmount.Add(taker.GetLaveAmount())
			maker.MatchAmount.Add(taker.GetLaveAmount())
			// 由於最後一個還有為搓合的數量，所以移除掉
			makers = makers[:len(makers)-1]
		}
	}

	return &MatchResult{
		Taker:        taker,
		FinishMakers: makers,
		Trades:       trades,
	}
}
