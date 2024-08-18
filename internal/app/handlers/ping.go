package handlers

import (
	"context"
	"net/http"
)

func PingHandler(ctx context.Context) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if (req.Method != http.MethodPost) &&
			(req.Method != http.MethodGet) &&
			(req.Method != http.MethodHead) {
			http.Error(res, "ping only accepts POST, GET and HEAD", http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Header().Add("Content-Type", "text/plain")
		res.Write([]byte("pong"))
	}

}
