package services_test

import (
	"RBKproject4/internal/services"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListTemplates_Success(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{"report.docx", "invoice.xlsx", "notes.txt"}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("dummy"), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", f, err)
		}
	}

	// Create a nested directory to test directory skipping
	if err := os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	svc := services.NewDocumentService(nil, nil, "", tmpDir, "", nil)

	result, err := svc.ListTemplates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"report":  "docx",
		"invoice": "xlsx",
		"notes":   "txt",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d templates, got %d", len(expected), len(result))
	}

	for _, tmpl := range result {
		if expFormat, ok := expected[tmpl.Name]; !ok {
			t.Errorf("unexpected template name: %s", tmpl.Name)
		} else if tmpl.Format != expFormat {
			t.Errorf("for template %s, expected format %s, got %s", tmpl.Name, expFormat, tmpl.Format)
		}
	}
}

func TestListTemplates_DirNotExist(t *testing.T) {
	// Point service to a non-existent dir
	svc := services.NewDocumentService(nil, nil, "", "nonexistent_dir", "", nil)

	_, err := svc.ListTemplates(context.Background())
	if err == nil {
		t.Fatal("expected an error when directory does not exist, got nil")
	}

	if !strings.Contains(err.Error(), "error listing templates") {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestListTemplates_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	svc := services.NewDocumentService(nil, nil, "", tmpDir, "", nil)

	result, err := svc.ListTemplates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty result, got %d templates", len(result))
	}
}

func TestListTemplates_FilesWithoutExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files without extensions
	files := []string{"README", "Makefile", "config"}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("dummy"), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", f, err)
		}
	}

	svc := services.NewDocumentService(nil, nil, "", tmpDir, "", nil)

	result, err := svc.ListTemplates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify behavior - service handles files without extensions
	for _, tmpl := range result {
		t.Logf("Found template: %s with format: %s", tmpl.Name, tmpl.Format)
	}
}
