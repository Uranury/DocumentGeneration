package renderers

import (
	"github.com/flosch/pongo2/v6"
)

type Pongo2Renderer struct {
	templateDir string
}

func NewPongo2Renderer(templateDir string) *Pongo2Renderer {
	return &Pongo2Renderer{templateDir: templateDir}
}

func (r *Pongo2Renderer) Render(templateName string, data map[string]interface{}) (string, error) {
	tpl, err := pongo2.FromFile(r.templateDir + "/" + templateName + ".html")
	if err != nil {
		return "", err
	}
	return tpl.Execute(data)
}
