package handler

import (
	"net/http"

	"github.com/Vla8islav/gophemart/internal/domain"
	"go.uber.org/zap"
)

type Handler struct {
	service domain.GophemartService
	logger  *zap.Logger
}

func NewHandler(service domain.GophemartService, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) writeUnauthorised(w http.ResponseWriter, msg string) {
	h.logger.Error("unauthorised", zap.String("msg", msg))
	http.Error(w, msg, http.StatusUnauthorized)
}

func (h *Handler) writeAlreadyExists(w http.ResponseWriter, msg string) {
	h.logger.Error("already exists", zap.String("msg", msg))
	http.Error(w, msg, http.StatusConflict)
}

func (h *Handler) writeInternalServerError(w http.ResponseWriter, msg string) {
	h.logger.Error("internal server error", zap.String("msg", msg))
	http.Error(w, msg, http.StatusInternalServerError)
}

func (h *Handler) writeBadRequest(w http.ResponseWriter, msg string) {
	h.logger.Error("bad request", zap.String("msg", msg))
	http.Error(w, msg, http.StatusBadRequest)
}

func (h *Handler) writeUnprocessableEntity(w http.ResponseWriter, msg string) {
	h.logger.Error("incorrect request checksum", zap.String("msg", msg))
	http.Error(w, msg, http.StatusUnprocessableEntity)
}

func (h *Handler) writeMethodNotAllowed(w http.ResponseWriter, msg string) {
	h.logger.Error("method not allowed: ", zap.String("msg", msg))
	http.Error(w, msg, http.StatusMethodNotAllowed)
}
