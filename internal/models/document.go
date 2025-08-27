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

func (d *Document) ContentType() string {
	switch d.Format {
	case FormatPDF:
		return "application/pdf"
	case FormatDOCX:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case FormatXLSX:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case FormatHTML:
		return "text/html"
	default:
		return "application/octet-stream"
	}
}
