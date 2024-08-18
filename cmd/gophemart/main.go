package main

import (
	"github.com/Vla8islav/gophemart/internal/app/handlers"
	"github.com/Vla8islav/gophemart/internal/app/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	ctx, cancel := helpers.GetDefaultContext()
	defer cancel()

	r := mux.NewRouter()

	r.HandleFunc("/ping/", handlers.PingHandler(ctx))

	err := http.ListenAndServe(helpers.ReadFlags().ServerAddress, r)
	if err != nil {

		panic(err)
	}
}
