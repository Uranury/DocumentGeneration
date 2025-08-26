package models

type DocumentFormat string

const (
	FormatDOCX DocumentFormat = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	FormatPDF  DocumentFormat = "application/pdf"
	FormatHTML DocumentFormat = "text/html"
	FormatXLSX DocumentFormat = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

type Document struct {
	Data     []byte
	Format   DocumentFormat
	Filename string
}
