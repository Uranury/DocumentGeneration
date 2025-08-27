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
	"os"
	"path/filepath"
)

type DocumentService struct {
	templateRenderer renderers.TemplateRenderer
	pythonURL        string
	templateDir      string
	gotenbergURL     string
	gotenbergPDFURL  string
	client           *http.Client
	logger           *slog.Logger
}

func NewDocumentService(logger *slog.Logger, templateRenderer renderers.TemplateRenderer, pythonURL, templateDir, gotenbergURL string, client *http.Client) *DocumentService {
	return &DocumentService{
		logger:           logger,
		pythonURL:        pythonURL,
		templateDir:      templateDir,
		templateRenderer: templateRenderer,
		gotenbergURL:     gotenbergURL,
		gotenbergPDFURL:  gotenbergURL + "/forms/chromium/convert/html",
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

func (s *DocumentService) renderWithPython(ctx context.Context, code, format string, data any) ([]byte, string, error) {
	route := "/docx/render" // For DOCX

	// Open template file
	templatePath := filepath.Join(s.templateDir, fmt.Sprintf("%s.%s", code, format))
	file, err := os.Open(templatePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open template: %w", err)
	}
	defer file.Close()

	jsonData, _ := json.Marshal(data)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("template", fmt.Sprintf("%s.%s", code, format))
	io.Copy(part, file)

	part, _ = writer.CreateFormFile("data", "data.json")
	part.Write(jsonData)

	writer.Close()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, s.pythonURL+route, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("python service returned status %d", resp.StatusCode)
	}

	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return dataBytes, string(models.FormatDOCX), nil
}

func (s *DocumentService) GenerateDOCX(ctx context.Context, req *models.RequestBody) (*models.Document, error) {
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data: %w", err)
	}

	dataBytes, _, err := s.renderWithPython(ctx, req.Code, "docx", dataMap)
	if err != nil {
		return nil, err
	}

	return &models.Document{
		Data:     dataBytes,
		Format:   models.FormatDOCX,
		Filename: fmt.Sprintf("%s.docx", req.Code),
	}, nil
}
