package run

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Vla8islav/gophemart/internal/accrual_client"
	"github.com/Vla8islav/gophemart/internal/config"
	"github.com/Vla8islav/gophemart/internal/domain"
	"github.com/Vla8islav/gophemart/internal/handler"
	"github.com/Vla8islav/gophemart/internal/middlewares"
	"github.com/Vla8islav/gophemart/internal/service"
	"go.uber.org/zap"
)

func Run(ctx context.Context, db domain.GophermartRepository, cfg *config.OptionsServer, logger *zap.Logger) error {

	gophermartAccrualClient := accrual_client.NewAccrualClient(cfg.AccrualAddress.Value, nil)

	srvApp := service.NewMetricsService(db, gophermartAccrualClient, cfg.AuthTokenSecret.Value)
	h := handler.NewHandler(srvApp, logger)
	r := handler.NewRouter(h, cfg)

	handlerWithMW := middlewares.ChainMiddlewares(
		r,
		middlewares.WithLogging(logger),
	)

	srv := &http.Server{
		Addr:         cfg.ServerAddress.Value,
		Handler:      handlerWithMW,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
