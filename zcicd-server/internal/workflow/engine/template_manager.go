package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// SystemTemplate represents a system preset build template.
type SystemTemplate struct {
	Name         string
	Language     string
	Framework    string
	Description  string
	TemplateFile string // relative path to YAML template
}

// TemplateManager manages build templates (loading from DB and rendering).
type TemplateManager struct {
	templateDir string
}

// NewTemplateManager creates a new TemplateManager with the given template directory.
func NewTemplateManager(templateDir string) *TemplateManager {
	return &TemplateManager{
		templateDir: templateDir,
	}
}

// GetSystemTemplates returns the list of system preset templates.
func (m *TemplateManager) GetSystemTemplates() []SystemTemplate {
	return []SystemTemplate{
		{
			Name:         "go-build",
			Language:     "go",
			Framework:    "",
			Description:  "Go application build with Kaniko image push",
			TemplateFile: "go_build.yaml",
		},
		{
			Name:         "java-maven",
			Language:     "java",
			Framework:    "maven",
			Description:  "Java Maven build with Kaniko image push",
			TemplateFile: "java_maven.yaml",
		},
		{
			Name:         "nodejs",
			Language:     "nodejs",
			Framework:    "",
			Description:  "Node.js build with npm/yarn/pnpm support and Kaniko image push",
			TemplateFile: "nodejs.yaml",
		},
	}
}

// LoadTemplate loads a template file and returns its content.
func (m *TemplateManager) LoadTemplate(templateFile string) (string, error) {
	path := filepath.Join(m.templateDir, templateFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load template %s: %w", templateFile, err)
	}
	return string(data), nil
}

// RenderTemplate renders a template with the given data.
func (m *TemplateManager) RenderTemplate(templateContent string, data interface{}) ([]byte, error) {
	tmpl, err := template.New("tekton").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	return buf.Bytes(), nil
}
