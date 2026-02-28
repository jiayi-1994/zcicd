package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const pipelineRunTemplate = `apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: wf-{{ .WorkflowID }}-run-{{ .RunNumber }}
  namespace: {{ .Namespace }}
  labels:
    app.kubernetes.io/managed-by: zcicd
    zcicd.io/workflow-id: "{{ .WorkflowID }}"
    zcicd.io/run-id: "{{ .RunID }}"
    zcicd.io/project-id: "{{ .ProjectID }}"
spec:
  pipelineSpec:
    tasks:
{{- range $si, $stage := .Stages }}
    - name: {{ $stage.Name }}
{{- if gt $si 0 }}
      runAfter:
      - {{ (index $.Stages (add $si -1)).Name }}
{{- end }}
      taskSpec:
        steps:
{{- range $job := $stage.Jobs }}
        - name: {{ $job.Name }}
          image: alpine:latest
          script: |
            echo "Executing job {{ $job.Name }} (type={{ $job.JobType }})"
{{- if gt $job.Timeout 0 }}
          timeout: {{ $job.Timeout }}s
{{- end }}
{{- end }}
{{- end }}
{{- if .Params }}
    params:
{{- range $key, $val := .Params }}
    - name: {{ $key }}
      type: string
      default: "{{ $val }}"
{{- end }}
{{- end }}
  workspaces:
  - name: shared-workspace
    volumeClaimTemplate:
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
`

const taskRunTemplate = `apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  name: build-{{ .BuildConfigID }}-run-{{ .RunNumber }}
  namespace: {{ .Namespace }}
  labels:
    app.kubernetes.io/managed-by: zcicd
    zcicd.io/build-config-id: "{{ .BuildConfigID }}"
    zcicd.io/run-id: "{{ .RunID }}"
    zcicd.io/project-id: "{{ .ProjectID }}"
    zcicd.io/service-name: "{{ .ServiceName }}"
spec:
  taskSpec:
    params:
    - name: repo_url
      type: string
    - name: branch
      type: string
    - name: commit_sha
      type: string
    - name: image_repo
      type: string
    - name: image_tag
      type: string
    steps:
    - name: git-clone
      image: alpine/git:latest
      script: |
        git clone --branch $(params.branch) --single-branch $(params.repo_url) /workspace/source
        cd /workspace/source
        git checkout $(params.commit_sha)
    - name: build
      image: alpine:latest
      workingDir: /workspace/source
      script: |
        {{ .BuildScript }}
    - name: docker-build
      image: gcr.io/kaniko-project/executor:latest
      args:
      - --dockerfile={{ .DockerfilePath }}
      - --context=/workspace/source/{{ .DockerContext }}
      - --destination=$(params.image_repo):$(params.image_tag)
{{- if .CacheEnabled }}
      - --cache=true
{{- end }}
    workspaces:
    - name: source
      description: Source code workspace
    - name: docker-config
      description: Docker config for registry auth
  params:
  - name: repo_url
    value: "{{ .RepoURL }}"
  - name: branch
    value: "{{ .Branch }}"
  - name: commit_sha
    value: "{{ .CommitSHA }}"
  - name: image_repo
    value: "{{ .ImageRepo }}"
  - name: image_tag
    value: "{{ .ImageTag }}"
  workspaces:
  - name: source
    emptyDir: {}
  - name: docker-config
    secret:
      secretName: docker-registry-credentials
`

// TemplateEngine renders platform workflow/build models into Tekton YAML.
type TemplateEngine struct {
	templateDir string
}

// NewTemplateEngine creates a new TemplateEngine with the given template directory.
func NewTemplateEngine(templateDir string) *TemplateEngine {
	return &TemplateEngine{templateDir: templateDir}
}

var templateFuncs = template.FuncMap{
	"add": func(a, b int) int { return a + b },
}

// RenderPipelineRun renders a WorkflowModel into Tekton PipelineRun YAML.
func (e *TemplateEngine) RenderPipelineRun(model *WorkflowModel) ([]byte, error) {
	tmpl, err := template.New("pipelinerun").Funcs(templateFuncs).Parse(pipelineRunTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pipeline run template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, model); err != nil {
		return nil, fmt.Errorf("failed to render pipeline run: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderTaskRun renders a BuildModel into Tekton TaskRun YAML.
func (e *TemplateEngine) RenderTaskRun(model *BuildModel) ([]byte, error) {
	tmpl, err := template.New("taskrun").Funcs(templateFuncs).Parse(taskRunTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task run template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, model); err != nil {
		return nil, fmt.Errorf("failed to render task run: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderFromTemplate renders a custom template file with the given data.
func (e *TemplateEngine) RenderFromTemplate(templateName string, data interface{}) ([]byte, error) {
	tmplPath := filepath.Join(e.templateDir, templateName)
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Funcs(templateFuncs).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to render template %s: %w", templateName, err)
	}
	return buf.Bytes(), nil
}
