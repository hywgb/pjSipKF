package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/hywgb/pjSipKF/control-plane/internal/config"
	"github.com/hywgb/pjSipKF/control-plane/internal/httpserver"
	"github.com/hywgb/pjSipKF/control-plane/internal/logging"
)

func main() {
	cfg := config.Load()
	logger := logging.NewLogger(cfg)
	defer logger.Sync()

	router := chi.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	// Register business routes
	httpserver.RegisterRoutes(router, logger, cfg)

	srv := &http.Server{
		Addr:              cfg.HTTPListen,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("control-plane HTTP server starting", zap.String("addr", cfg.HTTPListen))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	} else {
		logger.Info("server shut down cleanly")
	}
}