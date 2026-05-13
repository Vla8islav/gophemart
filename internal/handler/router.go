package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	r.Get("/api/ping", h.DBPing)

	r.Post("/api/user/register", h.UserRegisterHandler)

	r.Post("/api/user/login", h.UserLoginHandler)

	r.Post("/api/user/orders", h.DummyHandler)

	r.Get("/api/user/balance", h.DummyHandler)

	r.Post("/api/user/balance/withdraw", h.DummyHandler)

	r.Get("/api/user/withdrawals", h.DummyHandler)

	return r
}
