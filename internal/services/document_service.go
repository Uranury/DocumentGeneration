package services

import (
	"RBKproject4/internal/models"
	"RBKproject4/internal/renderers"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

type DocumentService struct {
	templateRenderer renderers.TemplateRenderer
	gotenbergURL     string
	gotenbergPDFURL  string
	gotenbergDocxURL string
	client           *http.Client
	logger           *slog.Logger
}

func NewDocumentService(logger *slog.Logger, templateRenderer renderers.TemplateRenderer, gotenbergURL string, client *http.Client) *DocumentService {
	return &DocumentService{
		logger:           logger,
		templateRenderer: templateRenderer,
		gotenbergURL:     gotenbergURL,
		gotenbergPDFURL:  gotenbergURL + "/forms/chromium/convert/html",
		gotenbergDocxURL: gotenbergURL + "/forms/libreoffice/convert?target=docx",
		client:           client,
	}
}

func toMap(data any) (map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *DocumentService) GenerateDocx(ctx context.Context, req *models.RequestBody) (*models.Document, error) {
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	renderedHTML, err := s.templateRenderer.Render(req.Code, dataMap)
	if err != nil {
		return nil, fmt.Errorf("error rendering html: %w", err)
	}

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("files", "input.html")
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	if _, err := part.Write([]byte(renderedHTML)); err != nil {
		return nil, fmt.Errorf("error writing to form file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing multipart writer: %w", err)
	}

	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.gotenbergDocxURL, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	newReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(newReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("failed to close response body")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gotenberg returned status %d, body: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return &models.Document{
		Data:     data,
		Format:   models.FormatDOCX,
		Filename: "document.docx",
	}, nil
}

func (s *DocumentService) GeneratePDF(ctx context.Context, req *models.RequestBody) (*models.Document, error) {
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	renderedHTML, err := s.templateRenderer.Render(req.Code, dataMap)
	if err != nil {
		return nil, fmt.Errorf("error rendering html: %w", err)
	}

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	_, err = part.Write([]byte(renderedHTML))
	if err != nil {
		return nil, fmt.Errorf("error writing to form file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing form file: %w", err)
	}

	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.gotenbergPDFURL, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	newReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(newReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("gotenberg returned status %d", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("failed to close response body")
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return &models.Document{
		Data:     data,
		Format:   models.FormatPDF,
		Filename: "document.pdf",
	}, nil
}

func (s *DocumentService) GenerateHTML(_ context.Context, req *models.RequestBody) (*models.Document, error) {
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	renderedHTML, err := s.templateRenderer.Render(req.Code, dataMap)
	if err != nil {
		return nil, fmt.Errorf("error rendering html: %w", err)
	}

	return &models.Document{
		Data:     []byte(renderedHTML),
		Format:   models.FormatHTML,
		Filename: "document.html",
	}, nil
}
