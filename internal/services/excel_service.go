package services

import (
	"RBKproject4/internal/models"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

func (s *DocumentService) GenerateXLSX(_ context.Context, req *models.RequestBody) (*models.Document, error) {
	// Convert data to map for easier processing
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	// Load the Excel template based on the code
	templatePath := s.getXLSXTemplatePath(req.Code)
	if templatePath == "" {
		return nil, fmt.Errorf("no XLSX template found for code: %s", req.Code)
	}

	// Load the template file
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error opening XLSX template: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			s.logger.Warn("Failed to close Excel template file")
		}
	}()

	// Get the active sheet (assuming single sheet templates)
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no active sheet found in template")
	}

	// Process the template
	err = s.processXLSXTemplate(f, sheetName, dataMap)
	if err != nil {
		return nil, fmt.Errorf("error processing XLSX template: %w", err)
	}

	// Save to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error writing Excel file to buffer: %w", err)
	}

	return &models.Document{
		Data:     buf.Bytes(),
		Format:   models.FormatXLSX,
		Filename: s.getOutputFilename(req.Code),
	}, nil
}

// Get template path based on code
func (s *DocumentService) getXLSXTemplatePath(code string) string {
	// You can configure this path or make it configurable
	basePath := "templates/"
	switch code {
	case "CARD_STATEMENT":
		return basePath + "CARD_STATEMENT_EXCEL.xlsx"
	case "EXTENDED_STATEMENT_EXCEL":
		return basePath + "EXTENDED_STATEMENT.xlsx"
	default:
		return ""
	}
}

// Get output filename based on code
func (s *DocumentService) getOutputFilename(code string) string {
	switch code {
	case "CARD_STATEMENT":
		return "card_statement.xlsx"
	case "EXTENDED_STATEMENT_EXCEL":
		return "extended_statement.xlsx"
	default:
		return "document.xlsx"
	}
}

// Process XLSX template with data
func (s *DocumentService) processXLSXTemplate(f *excelize.File, sheetName string, dataMap map[string]interface{}) error {
	// First pass: replace simple placeholders
	err := s.replacePlaceholders(f, sheetName, dataMap)
	if err != nil {
		return fmt.Errorf("error replacing placeholders: %w", err)
	}

	// Second pass: handle table data (arrays)
	err = s.processTableData(f, sheetName, dataMap)
	if err != nil {
		return fmt.Errorf("error processing table data: %w", err)
	}

	return nil
}

// Replace simple {{placeholder}} values
func (s *DocumentService) replacePlaceholders(f *excelize.File, sheetName string, dataMap map[string]interface{}) error {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}

	placeholderRegex := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.]+)\s*\}\}`)

	for rowIndex, row := range rows {
		for colIndex, cellValue := range row {
			if cellValue == "" {
				continue
			}

			// Find and replace all placeholders in this cell
			newValue := placeholderRegex.ReplaceAllStringFunc(cellValue, func(match string) string {
				// Extract placeholder name
				placeholder := placeholderRegex.FindStringSubmatch(match)
				if len(placeholder) < 2 {
					return match
				}

				key := placeholder[1]

				// Handle nested keys like "client.name"
				value := s.getNestedValue(dataMap, key)
				if value != nil {
					return fmt.Sprintf("%v", value)
				}

				return match // Keep original if not found
			})

			// Update cell if changed
			if newValue != cellValue {
				cellAddr := s.getCellAddress(colIndex, rowIndex+1) // Excel is 1-indexed
				if err := f.SetCellValue(sheetName, cellAddr, newValue); err != nil {
					return fmt.Errorf("error setting cell value: %w", err)
				}
			}
		}
	}

	return nil
}

// Process table data (arrays) with row duplication
func (s *DocumentService) processTableData(f *excelize.File, sheetName string, dataMap map[string]interface{}) error {
	// Find table placeholders like {{table1.field}}
	tablePlaceholders := s.findTablePlaceholders(f, sheetName)

	for tableName, placeholders := range tablePlaceholders {
		tableData, exists := dataMap[tableName]
		if !exists {
			continue
		}

		// Convert to slice of maps
		tableSlice, ok := s.convertToTableSlice(tableData)
		if !ok {
			s.logger.Warn(fmt.Sprintf("Table data for '%s' is not an array", tableName))
			continue
		}

		if len(tableSlice) == 0 {
			continue
		}

		// Process this table
		err := s.processTable(f, sheetName, tableName, placeholders, tableSlice)
		if err != nil {
			return fmt.Errorf("error processing table %s: %w", tableName, err)
		}
	}

	return nil
}

// Find table placeholders like {{tableName.field}}
func (s *DocumentService) findTablePlaceholders(f *excelize.File, sheetName string) map[string]map[string]CellLocation {
	result := make(map[string]map[string]CellLocation)
	placeholderRegex := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\s*\}\}`)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return result
	}

	for rowIndex, row := range rows {
		for colIndex, cellValue := range row {
			if cellValue == "" {
				continue
			}

			matches := placeholderRegex.FindStringSubmatch(cellValue)
			if len(matches) >= 3 {
				tableName := matches[1]
				fieldName := matches[2]

				if result[tableName] == nil {
					result[tableName] = make(map[string]CellLocation)
				}

				result[tableName][fieldName] = CellLocation{
					Row: rowIndex + 1, // Excel is 1-indexed
					Col: colIndex,
				}
			}
		}
	}

	return result
}

type CellLocation struct {
	Row int
	Col int
}

// Process a single table by duplicating template rows
func (s *DocumentService) processTable(f *excelize.File, sheetName, _ string, placeholders map[string]CellLocation, tableData []map[string]interface{}) error {
	if len(placeholders) == 0 || len(tableData) == 0 {
		return nil
	}

	// Find the template row (minimum row number with placeholders)
	templateRow := int(^uint(0) >> 1) // Max int
	for _, loc := range placeholders {
		if loc.Row < templateRow {
			templateRow = loc.Row
		}
	}

	// Insert additional rows if needed
	rowsToInsert := len(tableData) - 1
	if rowsToInsert > 0 {
		err := f.InsertRows(sheetName, templateRow+1, rowsToInsert)
		if err != nil {
			return fmt.Errorf("error inserting rows: %w", err)
		}
	}

	// Fill data into all rows
	for rowIdx, rowData := range tableData {
		currentRow := templateRow + rowIdx

		// Copy formatting from template row to new rows
		if rowIdx > 0 {
			err := s.copyRowFormatting(f, sheetName, templateRow, currentRow, placeholders)
			if err != nil {
				s.logger.Warn(fmt.Sprintf("Failed to copy formatting for row %d: %v", currentRow, err))
			}
		}

		// Fill data
		for fieldName, location := range placeholders {
			if value, exists := rowData[fieldName]; exists {
				cellAddr := s.getCellAddress(location.Col, currentRow)
				if err := f.SetCellValue(sheetName, cellAddr, value); err != nil {
					s.logger.Warn(fmt.Sprintf("Failed to set cell value for row %d: %v", currentRow, err))
				}
			}
		}
	}

	return nil
}

// Copy formatting from source row to target row
func (s *DocumentService) copyRowFormatting(f *excelize.File, sheetName string, sourceRow, targetRow int, placeholders map[string]CellLocation) error {
	for _, location := range placeholders {
		sourceAddr := s.getCellAddress(location.Col, sourceRow)
		targetAddr := s.getCellAddress(location.Col, targetRow)

		// Get source cell style
		styleID, err := f.GetCellStyle(sheetName, sourceAddr)
		if err != nil {
			continue // Skip if can't get style
		}

		// Apply style to target cell
		err = f.SetCellStyle(sheetName, targetAddr, targetAddr, styleID)
		if err != nil {
			continue // Skip if can't set style
		}
	}
	return nil
}

// Helper: get nested value from map using dot notation
func (s *DocumentService) getNestedValue(data map[string]interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	current := data

	for i, k := range keys {
		if i == len(keys)-1 {
			// Last key
			return current[k]
		}

		// Navigate deeper
		if next, ok := current[k].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

// Helper: convert interface{} to slice of maps
func (s *DocumentService) convertToTableSlice(data interface{}) ([]map[string]interface{}, bool) {
	switch v := data.(type) {
	case []interface{}:
		result := make([]map[string]interface{}, 0, len(v))
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				result = append(result, itemMap)
			} else {
				// Try to convert to map
				if converted, err := toMap(item); err == nil {
					result = append(result, converted)
				}
			}
		}
		return result, true
	case []map[string]interface{}:
		return v, true
	default:
		return nil, false
	}
}

// Helper: get cell address like A1, B2, etc.
func (s *DocumentService) getCellAddress(col, row int) string {
	return fmt.Sprintf("%s%d", s.getColumnName(col), row)
}

// Helper: get column name (A, B, C, ..., Z, AA, AB, etc.)
func (s *DocumentService) getColumnName(index int) string {
	name := ""
	for index >= 0 {
		name = string(rune('A'+index%26)) + name
		index = index/26 - 1
	}
	return name
}
