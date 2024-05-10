package handler

import (
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
	orderBook *model.OrderBook
	db        *mysql.Database
}

func NewHandler(db *mysql.Database) *Handler {
	h := new(Handler)
	h.orderBook = model.NewOrderBook(model.QueueTypePriceTree)
	h.db = db
	return h
}

type ReqCreateOrder struct {
	Price  *decimal.Decimal `json:"price"`
	Amount *decimal.Decimal `json:"amount"`
	Side   *model.OrderSide `json:"side"`
	Type   *model.OrderType `json:"type"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	reqBody := new(ReqCreateOrder)
	err := c.Bind(&reqBody)

	if err != nil {
		c.Error(err)
		return
	}

	// 需要先檢查訂單簿裡有沒有 避免資料被覆蓋 或者重複 - 可以放進外部快取
	order := &model.Order{
		Price:  *reqBody.Price,
		Amount: *reqBody.Amount,
		Side:   *reqBody.Side,
		Type:   *reqBody.Type,
	}

	if err = h.db.CreateOrder(order); err != nil {
		return
	}

	h.orderBook.Match(order)
	c.JSON(http.StatusCreated, order)
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	idStr := c.Param("order_id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil || id <= 0 {
		log.Warn().Err(err).Send()
		c.Status(http.StatusBadRequest)
		return
	}

	order, err := h.db.GetOrder(uint64(id))
	if err != nil {
		log.Warn().Err(err).Send()
		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusBadRequest)
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}
	err = h.db.DeleteOrder(uint64(id))
	if err != nil {
		log.Warn().Err(err).Send()
		c.Status(http.StatusInternalServerError)
		return
	}
	h.orderBook.CancelOrder(&order)
	c.Status(http.StatusNoContent)
}
