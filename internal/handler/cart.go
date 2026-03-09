package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ecommerce/internal/middleware"
	"ecommerce/internal/service"
)

func (h *Handler) GetCart(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)

	items, err := h.svc.Cart.Get(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *Handler) AddToCart(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)

	var in service.AddToCartInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Cart.Add(c.Request.Context(), userID, in); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) RemoveFromCart(c *gin.Context) {
	userID := c.GetInt64(middleware.CtxUserID)

	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	if err := h.svc.Cart.Remove(c.Request.Context(), userID, productID); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
