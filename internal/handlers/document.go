package handlers

import (
	"RBKproject4/internal/models"
	"RBKproject4/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DocumentHandler struct {
	svc *services.DocumentService
}

func NewDocumentHandler(svc *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

func (h *DocumentHandler) GenerateDocument(c *gin.Context) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var (
		document *models.Document
		err      error
	)

	switch req.Format {
	case "pdf":
		document, err = h.svc.GeneratePDF(c.Request.Context(), &req)
	case "docx":
		document, err = h.svc.GenerateDocx(c.Request.Context(), &req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", string(document.Format))
	c.Header("Content-Disposition", "attachment; filename="+document.Filename)
	c.Data(http.StatusOK, string(document.Format), document.Data)
}

func (h *DocumentHandler) GenerateHTML(c *gin.Context) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var (
		document *models.Document
		err      error
	)

	switch req.Format {
	case "pdf":
		document, err = h.svc.GeneratePDF(c.Request.Context(), &req)
	case "html":
		document, err = h.svc.GenerateHTML(c.Request.Context(), &req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", string(document.Format))
	c.Header("Content-Disposition", "attachment; filename="+document.Filename)
	c.Data(http.StatusOK, string(document.Format), document.Data)
}

func (h *DocumentHandler) GenerateXlsx(c *gin.Context) {
	var req models.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Format != "xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
	}

	document, err := h.svc.GenerateXLSX(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", string(document.Format))
	c.Header("Content-Disposition", "attachment; filename="+document.Filename)
	c.Data(http.StatusOK, string(document.Format), document.Data)
}
