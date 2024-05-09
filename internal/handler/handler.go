package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	h := new(Handler)
	return h
}

func (h *Handler) CreateOrder(ctx *gin.Context) {
	panic("not implement yet")

}

func (h *Handler) DeleteOrder(ctx *gin.Context) {
	panic("not implement yet")
}
