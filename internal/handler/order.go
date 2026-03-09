package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ecommerce/internal/middleware"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)

	order, err := h.svc.Order.CreateFromCart(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) ListOrders(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)

	orders, err := h.svc.Order.ListByUser(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *Handler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.svc.Order.GetByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}
