package server

import (
	"RBKproject4/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) setupRoutes() {
	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	api := s.Router.Group("/api/v1")
	docGeneration := api.Group(s.Cfg.ServiceContextURL)
	docGeneration.Use(middleware.AuthMiddleware(s.Cfg.StaticToken))

	docGeneration.POST("/generate-docx", s.DocumentHandler.GenerateDocument)
	docGeneration.POST("/generate-xlsx", s.DocumentHandler.GenerateXlsx)
	docGeneration.POST("/generate-html", s.DocumentHandler.GenerateHTML)
}
