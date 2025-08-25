package handlers

import (
	"RBKproject4/internal/models"
	"RBKproject4/internal/services"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DocumentHandler struct {
	svc *services.DocumentService
}

func NewDocumentHandler(svc *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func (h *DocumentHandler) generate(
	c *gin.Context,
	generateFunc func(ctx context.Context, req *models.RequestBody) ([]byte, services.DocumentFormat, string, error),
) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, format, filename, err := generateFunc(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", string(format))
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, string(format), doc)
}

func (h *DocumentHandler) GenerateDocx(c *gin.Context) {
	h.generate(c, h.svc.GenerateDocx)
}

func (h *DocumentHandler) GeneratePDF(c *gin.Context) {
	h.generate(c, h.svc.GeneratePDF)
}

func (h *DocumentHandler) GenerateHTML(c *gin.Context) {
	h.generate(c, h.svc.GenerateHTML)
}
