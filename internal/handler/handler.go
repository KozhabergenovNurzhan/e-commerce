package handler

import (
	"ecommerce/internal/pkg"
	"ecommerce/internal/service"
)

type Handler struct {
	svc *service.Services
}

func New(svc *service.Services) *Handler {
	return &Handler{svc: svc}
}

var respondError = pkg.RespondError
