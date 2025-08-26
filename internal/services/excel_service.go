package services

import (
	"RBKproject4/internal/models"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
	"regexp"
	"strconv"
	"strings"
)

// TODO: Change formatting, handle unhandled errors

func (s *DocumentService) GenerateXLSX(_ context.Context, req *models.RequestBody) (*models.Document, error) {
	// First, render HTML template like other methods
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	renderedHTML, err := s.templateRenderer.Render(req.Code, dataMap)
	if err != nil {
		return nil, fmt.Errorf("error rendering html: %w", err)
	}

	// Create new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			s.logger.Warn("Failed to close Excel file")
		}
	}()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(renderedHTML))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %w", err)
	}

	// Create main worksheet
	sheetName := "Document"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error creating sheet: %w", err)
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Convert HTML structure to Excel
	err = s.htmlToExcel(f, sheetName, doc)
	if err != nil {
		return nil, fmt.Errorf("error converting HTML to Excel: %w", err)
	}

	// Save to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("error writing Excel file to buffer: %w", err)
	}

	return &models.Document{
		Data:     buf.Bytes(),
		Format:   models.FormatXLSX,
		Filename: "document.xlsx",
	}, nil
}

// Generic method to convert any HTML structure to Excel
func (s *DocumentService) htmlToExcel(f *excelize.File, sheetName string, doc *goquery.Document) error {
	currentRow := 1

	// Process the body content sequentially
	doc.Find("body").Children().Each(func(i int, element *goquery.Selection) {
		currentRow = s.processElement(f, sheetName, element, currentRow, 0)
	})

	// Auto-adjust column widths
	s.autoSizeColumns(f, sheetName, currentRow)

	return nil
}

// Process any HTML element recursively
func (s *DocumentService) processElement(f *excelize.File, sheetName string, element *goquery.Selection, row, level int) int {
	tagName := goquery.NodeName(element)

	switch tagName {
	case "table":
		return s.processTable(f, sheetName, element, row)
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return s.processHeading(f, sheetName, element, row)
	case "p", "div", "span":
		return s.processTextElement(f, sheetName, element, row, level)
	case "ul", "ol":
		return s.processList(f, sheetName, element, row, level)
	default:
		// For other elements, process their children
		currentRow := row
		element.Children().Each(func(i int, child *goquery.Selection) {
			currentRow = s.processElement(f, sheetName, child, currentRow, level)
		})
		return currentRow
	}
}

// Process table elements
func (s *DocumentService) processTable(f *excelize.File, sheetName string, table *goquery.Selection, startRow int) int {
	currentRow := startRow

	// Process table headers
	table.Find("thead tr").Each(func(i int, row *goquery.Selection) {
		colIndex := 0
		row.Find("th, td").Each(func(j int, cell *goquery.Selection) {
			cellText := strings.TrimSpace(cell.Text())
			if cellText != "" {
				cellAddr := s.getCellAddress(colIndex, currentRow)
				f.SetCellValue(sheetName, cellAddr, cellText)
				s.setBoldStyle(f, sheetName, cellAddr)
			}
			colIndex++
		})
		currentRow++
	})

	// Process table body
	table.Find("tbody tr, tr:not(thead tr)").Each(func(i int, row *goquery.Selection) {
		// Skip if this row was already processed as header
		if table.Find("thead").Length() > 0 && i == 0 && row.Closest("thead").Length() > 0 {
			return
		}

		colIndex := 0
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			cellText := strings.TrimSpace(cell.Text())
			if cellText != "" {
				cellAddr := s.getCellAddress(colIndex, currentRow)

				// Try to parse as number
				if num, err := s.smartParseNumber(cellText); err == nil {
					f.SetCellValue(sheetName, cellAddr, num)
				} else {
					f.SetCellValue(sheetName, cellAddr, cellText)
				}
			}
			colIndex++
		})
		currentRow++
	})

	return currentRow + 1 // Add spacing after table
}

// Process heading elements
func (s *DocumentService) processHeading(f *excelize.File, sheetName string, heading *goquery.Selection, row int) int {
	text := strings.TrimSpace(heading.Text())
	if text == "" {
		return row
	}

	f.SetCellValue(sheetName, s.getCellAddress(0, row), text)
	s.setBoldStyle(f, sheetName, s.getCellAddress(0, row))

	return row + 2 // Add spacing after heading
}

// Process text elements (p, div, span)
func (s *DocumentService) processTextElement(f *excelize.File, sheetName string, element *goquery.Selection, row, level int) int {
	text := strings.TrimSpace(element.Text())
	if text == "" {
		return row
	}

	// Check if this looks like a key-value pair
	if strings.Contains(text, ":") && !strings.Contains(text, "\n") {
		parts := strings.SplitN(text, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			f.SetCellValue(sheetName, s.getCellAddress(level, row), key+":")
			s.setBoldStyle(f, sheetName, s.getCellAddress(level, row))

			// Try to parse value as number
			if num, err := s.smartParseNumber(value); err == nil {
				f.SetCellValue(sheetName, s.getCellAddress(level+1, row), num)
			} else {
				f.SetCellValue(sheetName, s.getCellAddress(level+1, row), value)
			}

			return row + 1
		}
	}

	// Regular text
	f.SetCellValue(sheetName, s.getCellAddress(level, row), text)
	return row + 1
}

// Process list elements
func (s *DocumentService) processList(f *excelize.File, sheetName string, list *goquery.Selection, row, level int) int {
	currentRow := row

	list.Find("li").Each(func(i int, item *goquery.Selection) {
		text := strings.TrimSpace(item.Text())
		if text != "" {
			prefix := "• "
			if goquery.NodeName(list) == "ol" {
				prefix = fmt.Sprintf("%d. ", i+1)
			}
			f.SetCellValue(sheetName, s.getCellAddress(level, currentRow), prefix+text)
			currentRow++
		}
	})

	return currentRow + 1
}

// Helper method to get cell address (A1, B2, etc.)
func (s *DocumentService) getCellAddress(col, row int) string {
	return fmt.Sprintf("%s%d", s.getColumnName(col), row)
}

// Helper method to get column name (A, B, C, ..., Z, AA, AB, etc.)
func (s *DocumentService) getColumnName(index int) string {
	name := ""
	for index >= 0 {
		name = string(rune('A'+index%26)) + name
		index = index/26 - 1
	}
	return name
}

// Smart number parsing that handles various formats
func (s *DocumentService) smartParseNumber(text string) (float64, error) {
	// Remove common non-numeric characters but keep decimal separators
	cleaned := regexp.MustCompile(`[^\d\.\-\,\+]`).ReplaceAllString(text, "")

	if cleaned == "" {
		return 0, fmt.Errorf("no numeric content")
	}

	// Handle comma as decimal separator (European format)
	if strings.Count(cleaned, ",") == 1 && strings.Count(cleaned, ".") == 0 {
		cleaned = strings.Replace(cleaned, ",", ".", 1)
	}

	// Handle thousand separators
	if strings.Count(cleaned, ",") > 1 || (strings.Contains(cleaned, ",") && strings.Contains(cleaned, ".")) {
		// Remove commas (thousand separators)
		cleaned = strings.Replace(cleaned, ",", "", -1)
	}

	return strconv.ParseFloat(cleaned, 64)
}

// Set bold style for cells
func (s *DocumentService) setBoldStyle(f *excelize.File, sheetName, cell string) {
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	f.SetCellStyle(sheetName, cell, cell, style)
}

// Auto-size columns based on content
func (s *DocumentService) autoSizeColumns(f *excelize.File, sheetName string, maxRow int) {
	// Get used range to determine how many columns to adjust
	maxCol := 10 // Default to first 10 columns, adjust as needed

	for col := 0; col < maxCol; col++ {
		colName := s.getColumnName(col)
		maxWidth := 8.0 // Minimum width

		// Sample a few rows to estimate width
		sampleRows := []int{1, maxRow / 4, maxRow / 2, 3 * maxRow / 4, maxRow}
		for _, row := range sampleRows {
			if row > 0 {
				cellValue, _ := f.GetCellValue(sheetName, fmt.Sprintf("%s%d", colName, row))
				if cellValue != "" {
					// Rough estimate: each character is about 1.2 units wide
					width := float64(len(cellValue)) * 1.2
					if width > maxWidth {
						maxWidth = width
					}
				}
			}
		}

		// Cap maximum width
		if maxWidth > 50 {
			maxWidth = 50
		}

		f.SetColWidth(sheetName, colName, colName, maxWidth)
	}
}
