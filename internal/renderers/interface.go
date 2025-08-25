package renderers

type TemplateRenderer interface {
	Render(templateName string, data map[string]interface{}) (string, error)
}
