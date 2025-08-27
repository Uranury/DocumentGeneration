package services

import (
	"RBKproject4/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

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

func (s *DocumentService) renderWithPython(ctx context.Context, code, format string, data any) ([]byte, string, error) {
	var route, contentType string
	switch format {
	case "docx":
		route = "/docx/render"
		contentType = string(models.FormatDOCX)
	case "xlsx":
		route = "/xlsx/render"
		contentType = string(models.FormatXLSX)
	default:
		return nil, "", fmt.Errorf("unsupported format %s", format)
	}

	templatePath := filepath.Join(s.templateDir, fmt.Sprintf("%s.%s", code, format))
	file, err := os.Open(templatePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open template: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Warn("failed to close template file", "err", err)
		}
	}()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal data: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// template file
	part, err := writer.CreateFormFile("template", fmt.Sprintf("%s.%s", code, format))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file for template: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, "", fmt.Errorf("failed to copy template: %w", err)
	}

	// data as a form field, not a file
	if err := writer.WriteField("data", string(jsonData)); err != nil {
		return nil, "", fmt.Errorf("failed to write data field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.pythonURL+route, body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}
	s.logger.Info("The request is being sent:", "endpoint", s.pythonURL+route)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("failed to close response body", "err", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("python service returned status %d", resp.StatusCode)
	}

	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response: %w", err)
	}

	return dataBytes, contentType, nil
}
