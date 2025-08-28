package renderers_test

import (
	"RBKproject4/internal/renderers"
	"os"
	"path/filepath"
	"testing"
)

func TestPongo2Renderer_Render(t *testing.T) {
	// make temp dir for template
	tmpDir := t.TempDir()

	// create a template file
	templateContent := "Hello {{ name }}!"
	templatePath := filepath.Join(tmpDir, "greet.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatal(err)
	}

	// init renderer
	r := renderers.NewPongo2Renderer(tmpDir)

	// run Render
	out, err := r.Render("greet", map[string]interface{}{"name": "World"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// assert result
	want := "Hello World!"
	if out != want {
		t.Errorf("Render() = %q, want %q", out, want)
	}
}

func TestPongo2Renderer_Render_MissingTemplate(t *testing.T) {
	r := renderers.NewPongo2Renderer("nonexistent")

	_, err := r.Render("missing", nil)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
