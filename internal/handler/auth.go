package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ecommerce/internal/service"
)

func (h *Handler) Register(c *gin.Context) {
	var in service.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Auth.Register(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "email": user.Email})
}

func (h *Handler) Login(c *gin.Context) {
	var in service.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Auth.Login(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
