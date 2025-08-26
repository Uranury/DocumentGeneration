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
	"os/exec"
	"path/filepath"
	"strings"
)

type DocumentService struct {
	templateRenderer renderers.TemplateRenderer
	libreOfficePath  string
	gotenbergURL     string
	client           *http.Client
	logger           *slog.Logger
}

func NewDocumentService(logger *slog.Logger, templateRenderer renderers.TemplateRenderer, libreOfficePath string, gotenbergURL string, client *http.Client) *DocumentService {
	return &DocumentService{logger: logger, templateRenderer: templateRenderer, libreOfficePath: libreOfficePath, gotenbergURL: gotenbergURL, client: client}
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

	tmpHTML, err := os.CreateTemp("", "*.html")
	if err != nil {
		return nil, fmt.Errorf("error creating temp HTML file: %w", err)
	}
	defer func() {
		if err := os.Remove(tmpHTML.Name()); err != nil {
			s.logger.Warn("error removing temp HTML file", tmpHTML.Name())
		}
	}()

	if err := os.WriteFile(tmpHTML.Name(), []byte(renderedHTML), 0644); err != nil {
		return nil, fmt.Errorf("error writing HTML file: %w", err)
	}

	// Create temporary output directory for LibreOffice
	tmpDir, err := os.MkdirTemp("", "libreoffice-out-*")
	if err != nil {
		return nil, fmt.Errorf("error creating temp dir: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			s.logger.Warn("error removing temp dir", tmpDir)
		}
	}()

	// Run LibreOffice headless to convert HTML -> DOCX
	cmd := exec.CommandContext(ctx, "soffice", "--headless", "--convert-to", "docx", tmpHTML.Name(), "--outdir", tmpDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("LibreOffice conversion failed: %w, stderr: %s", err, stderr.String())
	}

	// Check if LibreOffice produced any files
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("error reading output directory: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("LibreOffice did not produce a DOCX file, stderr: %s", stderr.String())
	}

	// Find the DOCX file (LibreOffice might add a prefix or use a different naming convention)
	var docxPath string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".docx") {
			docxPath = filepath.Join(tmpDir, file.Name())
			break
		}
	}

	if docxPath == "" {
		return nil, fmt.Errorf("no DOCX file found in output directory, files: %v", files)
	}

	data, err := os.ReadFile(docxPath)
	if err != nil {
		return nil, fmt.Errorf("error reading DOCX file: %w", err)
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

	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.gotenbergURL+"/forms/chromium/convert/html", buf)
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
