package service

import (
	"context"
	"fmt"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepo
	cartRepo    *repository.CartRepo
	productRepo *repository.ProductRepo
}

func NewOrderService(
	orderRepo *repository.OrderRepo,
	cartRepo *repository.CartRepo,
	productRepo *repository.ProductRepo,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateFromCart(ctx context.Context, userID int64) (*models.Order, error) {
	cartItems, err := s.cartRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("%w: cart is empty", models.ErrInvalidInput)
	}

	var total float64
	orderItems := make([]models.OrderItem, 0, len(cartItems))

	for _, ci := range cartItems {
		product, err := s.productRepo.GetByID(ctx, ci.ProductID)
		if err != nil {
			return nil, err
		}
		if err := s.productRepo.DecreaseStock(ctx, ci.ProductID, ci.Quantity); err != nil {
			return nil, err
		}

		total += product.Price * float64(ci.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			ProductID: ci.ProductID,
			Quantity:  ci.Quantity,
			Price:     product.Price,
		})
	}

	order := &models.Order{
		UserID: userID,
		Status: models.StatusPending,
		Total:  total,
		Items:  orderItems,
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	_ = s.cartRepo.Clear(ctx, userID)

	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) ListByUser(ctx context.Context, userID int64) ([]models.Order, error) {
	return s.orderRepo.ListByUser(ctx, userID)
}
