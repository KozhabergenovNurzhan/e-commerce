package service

import (
	"context"

	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService(repo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func (s *CategoryService) Create(ctx context.Context, in CreateCategoryInput) (*models.Category, error) {
	c := &models.Category{Name: in.Name}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) List(ctx context.Context) ([]models.Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) Update(ctx context.Context, id int64, in CreateCategoryInput) (*models.Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Name = in.Name
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
