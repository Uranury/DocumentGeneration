package server

import (
	"RBKproject4/pkg/config"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	Router     *gin.Engine
	HTTPServer *http.Server
	Cfg        *config.Config
	Logger     *slog.Logger
}

func NewServer(cfg *config.Config, logger *slog.Logger) *Server {
	router := gin.Default()

	httpServer := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		Router:     router,
		Cfg:        cfg,
		Logger:     logger,
		HTTPServer: httpServer,
	}
}

func (s *Server) Run() error {
	s.Logger.Info("Starting server at address...", s.Cfg.ListenAddr)
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Shutting down server...")
	done := make(chan error, 1)

	go func() {
		if err := s.HTTPServer.Shutdown(ctx); err != nil {
			s.Logger.Error("Failed to shutdown http server", "error", err)
			done <- err
			return
		}
		s.Logger.Info("Graceful shutdown completed")
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		s.Logger.Info("Graceful shutdown timed out")
		return ctx.Err()
	}
}
