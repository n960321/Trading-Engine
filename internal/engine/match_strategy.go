package engine

import (
	"Trading-Engine/internal/model"

	"github.com/shopspring/decimal"
)

type MatchStrategy interface {
	MatchOrder(buyQueue, sellQueue OrderQueueior, taker model.Order) model.MatchResult
}

// 限價單撮合
type LimitStrategy struct{}

func (l LimitStrategy) MatchOrder(buyQueue, sellQueue OrderQueueior, taker model.Order) model.MatchResult {
	var matchResult model.MatchResult
	var makers []model.Order
	switch taker.Side {
	case model.OrderSideBuy:
		makers = sellQueue.GetOrdersBetweenPirceWithAmount(taker.Price, taker.GetLaveAmount())
	case model.OrderSideSell:
		makers = buyQueue.GetOrdersBetweenPirceWithAmount(taker.Price, taker.GetLaveAmount())
	}
	if len(makers) == 0 {
		if taker.Side == model.OrderSideBuy {
			buyQueue.AddOrder(taker)
		} else {
			sellQueue.AddOrder(taker)
		}
		matchResult.Taker = taker
		return matchResult
	}

	trades := make([]model.Trade, 0)
	for i, maker := range makers {
		// 當 taker 剩餘Amount 大於等於 maker的LaveAmount
		if taker.GetLaveAmount().Compare(maker.GetLaveAmount()) >= 0 {
			trades = append(trades, model.Trade{
				TakerID: taker.ID,
				MakerID: maker.ID,
				Amount:  maker.GetLaveAmount(),
				Price:   maker.Price,
			})

			taker.MatchAmount = taker.MatchAmount.Add(maker.GetLaveAmount())
			makers[i].MatchAmount = makers[i].MatchAmount.Add(makers[i].GetLaveAmount())
			if makers[i].GetLaveAmount().Equal(decimal.Zero) {
				makers[i].Compeleted = true
			}

		} else {
			// 當 taker 剩餘Amount 小於 maker的LaveAmount
			trades = append(trades, model.Trade{
				TakerID: taker.ID,
				MakerID: maker.ID,
				Amount:  taker.GetLaveAmount(),
				Price:   maker.Price,
			})

			makers[i].MatchAmount = makers[i].MatchAmount.Add(taker.GetLaveAmount())
			taker.MatchAmount = taker.MatchAmount.Add(taker.GetLaveAmount())
			matchResult.IncompleteMaker = makers[i]
			// 由於最後一個還有為搓合的數量，所以移除掉
			makers = makers[:len(makers)-1]
		}
	}

	// taker 剩餘Amount 是零要更動flag
	if taker.GetLaveAmount().Equal(decimal.Zero) {
		taker.Compeleted = true
	}

	matchResult.Taker = taker
	matchResult.Trades = trades
	matchResult.FinishMakers = makers
	return matchResult
}

// 市價單撮合
type MarketStrategy struct{}

func (l MarketStrategy) MatchOrder(buyQueue, sellQueue OrderQueueior, taker model.Order) model.MatchResult {
	var matchResult model.MatchResult
	var makers []model.Order
	switch taker.Side {
	case model.OrderSideBuy:
		makers = sellQueue.GetOrdersWithAmount(taker.GetLaveAmount())
	case model.OrderSideSell:
		makers = buyQueue.GetOrdersWithAmount(taker.GetLaveAmount())
	}

	matchResult.Taker = taker
	if len(makers) == 0 {
		matchResult.Taker.Compeleted = true
		return matchResult
	}
	trades := make([]model.Trade, 0)
	for i, maker := range makers {
		// 當 taker 剩餘Amount 大於等於 maker的LaveAmount
		if taker.GetLaveAmount().Compare(maker.GetLaveAmount()) >= 0 {
			trades = append(trades, model.Trade{
				TakerID: taker.ID,
				MakerID: maker.ID,
				Amount:  maker.GetLaveAmount(),
				Price:   maker.Price,
			})

			taker.MatchAmount = taker.MatchAmount.Add(maker.GetLaveAmount())
			makers[i].MatchAmount = makers[i].MatchAmount.Add(makers[i].GetLaveAmount())
			if makers[i].GetLaveAmount().Equal(decimal.Zero) {
				makers[i].Compeleted = true
			}

		} else {
			// 當 taker 剩餘Amount 小於 maker的LaveAmount
			trades = append(trades, model.Trade{
				TakerID: taker.ID,
				MakerID: maker.ID,
				Amount:  taker.GetLaveAmount(),
				Price:   maker.Price,
			})

			makers[i].MatchAmount = makers[i].MatchAmount.Add(taker.GetLaveAmount())
			taker.MatchAmount = taker.MatchAmount.Add(taker.GetLaveAmount())
			matchResult.IncompleteMaker = makers[i]
			// 由於最後一個還有為搓合的數量，所以移除掉
			makers = makers[:len(makers)-1]
		}
	}
	// 市價單 當下就結束
	taker.Compeleted = true
	matchResult.Taker = taker
	matchResult.Trades = trades
	matchResult.FinishMakers = makers
	return matchResult
}
