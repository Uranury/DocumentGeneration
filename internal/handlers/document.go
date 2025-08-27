package handlers

import (
	"RBKproject4/internal/models"
	"RBKproject4/internal/services"
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DocumentHandler struct {
	svc *services.DocumentService
}

func NewDocumentHandler(svc *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func streamDocument(c *gin.Context, doc *models.Document, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.DataFromReader(
		http.StatusOK,
		int64(len(doc.Data)),
		doc.ContentType(),
		bytes.NewReader(doc.Data),
		map[string]string{
			"Content-Disposition": "attachment; filename=" + doc.Filename,
		},
	)
}

func (h *DocumentHandler) GenerateHTML(c *gin.Context) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	var doc *models.Document
	var err error

	switch req.Format {
	case "html":
		doc, err = h.svc.GenerateHTML(ctx, &req)
	case "pdf":
		doc, err = h.svc.GeneratePDF(ctx, &req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	streamDocument(c, doc, err)
}

func (h *DocumentHandler) GenerateDocument(c *gin.Context) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	var doc *models.Document
	var err error

	switch req.Format {
	case "docx":
		doc, err = h.svc.GenerateDOCX(ctx, &req)
	case "pdf":
		doc, err = h.svc.GeneratePDF(ctx, &req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	streamDocument(c, doc, err)
}

func (h *DocumentHandler) GenerateXLSX(c *gin.Context) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Format != "xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
		return
	}

	doc, err := h.svc.GenerateXLSX(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	streamDocument(c, doc, err)
}
