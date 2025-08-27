package server

import (
	"RBKproject4/internal/handlers"
	"RBKproject4/internal/renderers"
	"RBKproject4/internal/services"
	"RBKproject4/pkg/config"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	Router          *gin.Engine
	HTTPServer      *http.Server
	HTTPClient      *http.Client
	DocumentHandler *handlers.DocumentHandler
	Cfg             *config.Config
	Logger          *slog.Logger
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

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	templateRenderer := renderers.NewPongo2Renderer(cfg.TemplateDir)
	newDocService := services.NewDocumentService(logger, templateRenderer, cfg.TemplateDir, cfg.PDFConverterURL, httpClient)
	newDocHandler := handlers.NewDocumentHandler(newDocService)

	server := &Server{
		Router:          router,
		Cfg:             cfg,
		Logger:          logger,
		HTTPServer:      httpServer,
		HTTPClient:      httpClient,
		DocumentHandler: newDocHandler,
	}

	server.setupRoutes()
	return server
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
