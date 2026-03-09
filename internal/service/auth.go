package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ecommerce/internal/auth"
	"ecommerce/internal/models"
	"ecommerce/internal/repository"
)

type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthService(userRepo *repository.UserRepo, secret string, ttl time.Duration) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: secret, jwtTTL: ttl}
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         models.RoleUser,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return "", models.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", models.ErrUnauthorized
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}
