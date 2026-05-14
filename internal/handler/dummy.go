package handler

import "net/http"

func (h *Handler) DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
