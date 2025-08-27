package server

import (
	"RBKproject4/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) setupRoutes() {
	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	api := s.Router.Group("/api/v1")
	docGeneration := api.Group(s.Cfg.ServiceContextURL)
	docGeneration.Use(middleware.AuthMiddleware(s.Cfg.StaticToken))

	docGeneration.POST("/generate-docx", s.DocumentHandler.GenerateDocument)
	docGeneration.POST("/generate-xlsx", s.DocumentHandler.GenerateXLSX)
	docGeneration.POST("/generate-html", s.DocumentHandler.GenerateHTML)
}
