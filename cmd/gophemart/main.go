package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	currentConfig := config.ReadFlagsServer(os.Args[1:])
	logger.Info("starting server ", zap.String("Server addr", currentConfig.ServerAddress.Value))

	srvImpl := &http.Server{Addr: currentConfig.ServerAddress.Value,
		Handler:      handlerWithMW,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	err = srvImpl.ListenAndServe()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
}
