package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"ecommerce/internal/handler"
	"ecommerce/internal/middleware"
	"ecommerce/internal/service"
)

func New(svc *service.Services, log *slog.Logger, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())

	h := handler.New(svc)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
		}

		api.GET("/products", h.ListProducts)
		api.GET("/products/:id", h.GetProduct)
		api.GET("/categories", h.ListCategories)

		protected := api.Group("")
		protected.Use(middleware.JWT(jwtSecret))
		{
			protected.GET("/cart", h.GetCart)
			protected.POST("/cart", h.AddToCart)
			protected.DELETE("/cart/:product_id", h.RemoveFromCart)

			protected.POST("/orders", h.CreateOrder)
			protected.GET("/orders", h.ListOrders)
			protected.GET("/orders/:id", h.GetOrder)

			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.POST("/products", h.CreateProduct)
				admin.PUT("/products/:id", h.UpdateProduct)
				admin.DELETE("/products/:id", h.DeleteProduct)

				admin.POST("/categories", h.CreateCategory)
				admin.PUT("/categories/:id", h.UpdateCategory)
				admin.DELETE("/categories/:id", h.DeleteCategory)
			}
		}
	}

	return r
}
