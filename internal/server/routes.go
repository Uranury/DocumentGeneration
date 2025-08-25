package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) setupRoutes() {
	s.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := s.Router.Group("/api/v1")
	docGeneration := api.Group(s.Cfg.ServiceContextURL)
}
