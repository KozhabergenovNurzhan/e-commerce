package service

import (
	"context"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

type CreateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryID  int64   `json:"category_id" binding:"required"`
}

func (s *ProductService) Create(ctx context.Context, in CreateProductInput) (*models.Product, error) {
	p := &models.Product{
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		Stock:       in.Stock,
		CategoryID:  in.CategoryID,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context, limit, offset int) ([]models.Product, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *ProductService) Update(ctx context.Context, p *models.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
