package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/api/ping", handler.DBPing)

	r.Post("/api/user/register", handler.DummyHandler)

	r.Post("/api/user/login", handler.DummyHandler)

	r.Post("/api/user/orders", handler.DummyHandler)

	r.Get("/api/user/balance", handler.DummyHandler)

	r.Post("/api/user/balance/withdraw", handler.DummyHandler)

	r.Get("/api/user/withdrawals", handler.DummyHandler)

	return r
}
