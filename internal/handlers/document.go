package handlers

import (
	"RBKproject4/internal/models"
	"RBKproject4/internal/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DocumentHandler struct {
	svc *services.DocumentService
}

func NewDocumentHandler(svc *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func (h *DocumentHandler) handleDocument(c *gin.Context, generate func(ctx context.Context, req *models.RequestBody) (*models.Document, error)) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document, err := generate(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", string(document.Format))
	c.Header("Content-Disposition", "attachment; filename="+document.Filename)
	c.Data(http.StatusOK, string(document.Format), document.Data)
}

func (h *DocumentHandler) GenerateDocument(c *gin.Context) {
	h.handleDocument(c, func(ctx context.Context, req *models.RequestBody) (*models.Document, error) {
		switch req.Format {
		case "pdf":
			return h.svc.GeneratePDF(ctx, req)
		case "docx":
			return h.svc.GenerateDocx(ctx, req)
		default:
			return nil, fmt.Errorf("unsupported format")
		}
	})
}

func (h *DocumentHandler) GenerateHTML(c *gin.Context) {
	h.handleDocument(c, func(ctx context.Context, req *models.RequestBody) (*models.Document, error) {
		switch req.Format {
		case "pdf":
			return h.svc.GeneratePDF(ctx, req)
		case "html":
			return h.svc.GenerateHTML(ctx, req)
		default:
			return nil, fmt.Errorf("unsupported format")
		}
	})
}

func (h *DocumentHandler) GenerateXlsx(c *gin.Context) {
	h.handleDocument(c, func(ctx context.Context, req *models.RequestBody) (*models.Document, error) {
		if req.Format != "xlsx" {
			return nil, fmt.Errorf("unsupported format")
		}
		return h.svc.GenerateXLSX(ctx, req)
	})
}
