package services

import (
	"RBKproject4/internal/models"
	"context"
	"fmt"
)

func (s *DocumentService) GenerateXLSX(ctx context.Context, req *models.RequestBody) (*models.Document, error) {
	dataMap, err := ToMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	dataBytes, _, err := s.renderWithPython(ctx, req.Code, "xlsx", dataMap)
	if err != nil {
		return nil, err
	}

	return &models.Document{
		Data:     dataBytes,
		Format:   models.FormatXLSX,
		Filename: fmt.Sprintf("%s.xlsx", req.Code),
	}, nil
}
