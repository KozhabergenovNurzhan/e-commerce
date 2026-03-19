package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ecommerce/internal/service"
)

func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.svc.Category.GetAll(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var in service.CreateCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.svc.Category.Create(c.Request.Context(), in)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var in service.CreateCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.svc.Category.Update(c.Request.Context(), id, in)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Category.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
