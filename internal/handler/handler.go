package handler

import (
	"github.com/Vla8islav/gophemart/internal/domain"
	"go.uber.org/zap"
)

type Handler struct {
	service domain.GophemartRepository
	logger  *zap.Logger
}

func NewHandler(service domain.GophemartService, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}
