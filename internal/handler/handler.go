package handler

import (
	"Trading-Engine/internal/engine"
	"Trading-Engine/internal/model"
	"Trading-Engine/internal/storage/mysql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Handler struct {
	orderBook *engine.OrderBook
	db        *mysql.Database
}

func NewHandler(db *mysql.Database) *Handler {
	h := new(Handler)
	h.orderBook = engine.NewOrderBook(engine.QueueTypePriceTree)
	h.db = db
	return h
}

type ReqCreateOrder struct {
	Price  decimal.Decimal `json:"price"`
	Amount decimal.Decimal `json:"amount"`
	Side   model.OrderSide `json:"side"`
	Type   model.OrderType `json:"type"`
}

func (h *Handler) GetOrder(c *gin.Context) {
	idStr := c.Param("order_id")

	if idStr == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Status(http.StatusBadRequest)
		log.Warn().Err(err).Msg("GetOrder parse uint failed")
		return
	}

	order, err := h.db.GetOrder(uint(id))

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusNotFound)
			return
		}
		log.Warn().Err(err).Msg("GetOrder from db failed")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	reqBody := new(ReqCreateOrder)
	err := c.Bind(&reqBody)

	if err != nil {
		c.Error(err)
		return
	}

	// 需要先檢查訂單簿裡有沒有 避免資料被覆蓋 或者重複 - 可以放進外部快取
	order := model.Order{
		Price:  reqBody.Price,
		Amount: reqBody.Amount,
		Side:   reqBody.Side,
		Type:   reqBody.Type,
	}

	if err = h.db.CreateOrder(nil, &order); err != nil {
		return
	}
	h.orderBook.Mux.Lock()
	defer h.orderBook.Mux.Unlock()
	matchResult := h.orderBook.Match(order)
	// h.orderBook.Match(order)

	// 更新進db
	{
		tx := h.db.GetTxBegin()
		// 	更新 taker 的狀態
		err = h.db.UpdateOrder(tx, &matchResult.Taker)
		if err != nil {
			tx.Rollback()
			c.Status(http.StatusInternalServerError)
			return
		}
		// 	更新 finishMarker
		for _, maker := range matchResult.FinishMakers {
			err = h.db.UpdateOrder(tx, &maker)
			if err != nil {
				tx.Rollback()
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		// 	更新 imcomplete marker
		if !matchResult.IncompleteMaker.CreatedAt.IsZero() {
			if err = h.db.UpdateOrder(tx, &matchResult.IncompleteMaker); err != nil {
				tx.Rollback()
				c.Status(http.StatusInternalServerError)
				return
			}
		}
		// 	新增 Trades
		for _, trade := range matchResult.Trades {
			if err := h.db.CreateTrade(tx, trade); err != nil {
				tx.Rollback()
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		tx.Commit()
	}
	// 更新orderbook的資料
	h.orderBook.UpdateFromMatchResult(*matchResult)

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("order_id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil || id <= 0 {
		log.Warn().Err(err).Send()
		c.Status(http.StatusBadRequest)
		return
	}

	h.orderBook.Mux.Lock()
	defer h.orderBook.Mux.Unlock()

	order, err := h.db.GetOrder(uint(id))
	if err != nil {
		log.Warn().Err(err).Send()
		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusBadRequest)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	if order.Compeleted {
		c.Status(http.StatusNotAcceptable)
		return
	}

	err = h.db.DeleteOrder(nil, uint(id))
	if err != nil {
		log.Warn().Err(err).Send()
		c.Status(http.StatusInternalServerError)
		return
	}
	h.orderBook.RemoveOrder(order)
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListTrades(c *gin.Context) {
	markIDStr := c.Query("makerID")
	takerIDStr := c.Query("takerID")

	if markIDStr == "" && takerIDStr == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	var makerID, takerID uint64
	var err error
	if markIDStr != "" {
		makerID, err = strconv.ParseUint(markIDStr, 10, 64)
		if err != nil {
			c.Status(http.StatusBadRequest)
			log.Warn().Err(err).Msg("GetOrder parse uint makerID failed")
			return
		}
	}

	if takerIDStr != "" {
		takerID, err = strconv.ParseUint(takerIDStr, 10, 64)
		if err != nil {
			c.Status(http.StatusBadRequest)
			log.Warn().Err(err).Msg("GetOrder parse uint takerID failed")
			return
		}
	}

	trades, err := h.db.ListTrades(uint(makerID), uint(takerID))

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusNotFound)
			return
		}
		log.Warn().Err(err).Msg("GetOrder from db failed")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, trades)
}
