package service

import (
	"context"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepo
	productRepo *repository.ProductRepo
}

func NewCartService(cartRepo *repository.CartRepo, productRepo *repository.ProductRepo) *CartService {
	return &CartService{cartRepo: cartRepo, productRepo: productRepo}
}

type AddToCartInput struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

func (s *CartService) Get(ctx context.Context, userID int64) ([]models.CartItem, error) {
	return s.cartRepo.GetByUser(ctx, userID)
}

func (s *CartService) Add(ctx context.Context, userID int64, in AddToCartInput) error {
	if _, err := s.productRepo.GetByID(ctx, in.ProductID); err != nil {
		return err
	}
	return s.cartRepo.Upsert(ctx, &models.CartItem{
		UserID:    userID,
		ProductID: in.ProductID,
		Quantity:  in.Quantity,
	})
}

func (s *CartService) Remove(ctx context.Context, userID, productID int64) error {
	return s.cartRepo.Remove(ctx, userID, productID)
}
