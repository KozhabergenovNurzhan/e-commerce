package service

import (
	"context"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

type CreateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	CategoryID  int64   `json:"category_id" binding:"required"`
}

type ListProductsInput struct {
	Search     string
	CategoryID int64
	Limit      int
	Offset     int
}

type ListProductsOutput struct {
	Data  []models.Product `json:"data"`
	Total int              `json:"total"`
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

func (s *ProductService) List(ctx context.Context, in ListProductsInput) (*ListProductsOutput, error) {
	if in.Limit <= 0 || in.Limit > 100 {
		in.Limit = 20
	}

	products, total, err := s.repo.List(ctx, repository.ProductFilter{
		Search:     in.Search,
		CategoryID: in.CategoryID,
		Limit:      in.Limit,
		Offset:     in.Offset,
	})
	if err != nil {
		return nil, err
	}

	return &ListProductsOutput{Data: products, Total: total}, nil
}

func (s *ProductService) Update(ctx context.Context, p *models.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
