package services

import (
	"RBKproject4/internal/models"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"baliance.com/gooxml/document"
)

func (s *DocumentService) GenerateDocx(_ context.Context, req *models.RequestBody) (*models.Document, error) {
	// Convert request data to map
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	// Construct template file path
	templatePath := filepath.Join(s.templateDir, req.Code+".docx")

	// Check if template file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template file not found: %s", templatePath)
	}

	// Open the DOCX template
	doc, err := document.Open(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error opening docx template: %w", err)
	}

	// Process placeholders in the document
	if err := s.replaceDocxPlaceholders(doc, dataMap); err != nil {
		return nil, fmt.Errorf("error replacing placeholders: %w", err)
	}

	// Create a buffer to write the modified document
	buf := &bytes.Buffer{}

	// Save the document to the buffer
	if err := doc.Save(buf); err != nil {
		return nil, fmt.Errorf("error saving docx document: %w", err)
	}

	return &models.Document{
		Data:     buf.Bytes(),
		Format:   models.FormatDOCX,
		Filename: "document.docx",
	}, nil
}

// replacePlaceholders replaces HTML-style placeholders in the document
func (s *DocumentService) replaceDocxPlaceholders(doc *document.Document, data map[string]interface{}) error {
	// Regular expression to match HTML-style placeholders like {{variable}} or {variable}
	placeholderRegex := regexp.MustCompile(`\{\{?\s*([^}]+)\s*\}?\}`)

	// Process all paragraphs in the document
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text := run.Text()
			if strings.Contains(text, "{") {
				// Find and replace all placeholders in this run
				newText := placeholderRegex.ReplaceAllStringFunc(text, func(match string) string {
					// Extract the variable name from the placeholder
					key := placeholderRegex.FindStringSubmatch(match)[1]
					key = strings.TrimSpace(key)

					// Look up the value in the data map
					if value, exists := data[key]; exists {
						return fmt.Sprintf("%v", value)
					}

					// If key not found, keep the original placeholder
					s.logger.Warn("placeholder not found in data", "key", key)
					return match
				})

				// Clear the run and add the new text
				run.Clear()
				run.AddText(newText)
			}
		}
	}

	// Process headers
	for _, hdr := range doc.Headers() {
		for _, para := range hdr.Paragraphs() {
			s.replacePlaceholdersInParagraph(para, data, placeholderRegex)
		}
	}

	// Process footers
	for _, ftr := range doc.Footers() {
		for _, para := range ftr.Paragraphs() {
			s.replacePlaceholdersInParagraph(para, data, placeholderRegex)
		}
	}

	// Process tables
	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					s.replacePlaceholdersInParagraph(para, data, placeholderRegex)
				}
			}
		}
	}

	return nil
}

// replacePlaceholdersInParagraph is a helper function to replace placeholders in a specific paragraph
func (s *DocumentService) replacePlaceholdersInParagraph(para document.Paragraph, data map[string]interface{}, placeholderRegex *regexp.Regexp) {
	for _, run := range para.Runs() {
		text := run.Text()
		if strings.Contains(text, "{") {
			newText := placeholderRegex.ReplaceAllStringFunc(text, func(match string) string {
				key := placeholderRegex.FindStringSubmatch(match)[1]
				key = strings.TrimSpace(key)

				if value, exists := data[key]; exists {
					return fmt.Sprintf("%v", value)
				}

				s.logger.Warn("placeholder not found in data", "key", key)
				return match
			})

			run.Clear()
			run.AddText(newText)
		}
	}
}

func (s *DocumentService) GenerateDocxFromHTML(_ context.Context, req *models.RequestBody) (*models.Document, error) {
	// First render the HTML template
	dataMap, err := toMap(req.Data)
	if err != nil {
		return nil, fmt.Errorf("error converting data to map: %w", err)
	}

	renderedHTML, err := s.templateRenderer.Render(req.Code, dataMap)
	if err != nil {
		return nil, fmt.Errorf("error rendering html: %w", err)
	}

	// Create a new DOCX document
	doc := document.New()

	// Simple HTML to DOCX conversion (basic implementation)
	// This is a simplified approach - for more complex HTML you might want to use a more sophisticated parser

	// Remove HTML tags and add content to document
	cleanText := s.stripHTMLTags(renderedHTML)
	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText(cleanText)

	// Save to buffer
	buf := &bytes.Buffer{}
	if err := doc.Save(buf); err != nil {
		return nil, fmt.Errorf("error saving docx document: %w", err)
	}

	return &models.Document{
		Data:     buf.Bytes(),
		Format:   models.FormatDOCX,
		Filename: "document.docx",
	}, nil
}

// stripHTMLTags removes HTML tags from text (basic implementation)
func (s *DocumentService) stripHTMLTags(html string) string {
	// Simple regex to remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(html, "")
}
