package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
)

// Tekton GVR definitions.
var (
	PipelineRunGVR = schema.GroupVersionResource{Group: "tekton.dev", Version: "v1", Resource: "pipelineruns"}
	TaskRunGVR     = schema.GroupVersionResource{Group: "tekton.dev", Version: "v1", Resource: "taskruns"}
	PipelineGVR    = schema.GroupVersionResource{Group: "tekton.dev", Version: "v1", Resource: "pipelines"}
	TaskGVR        = schema.GroupVersionResource{Group: "tekton.dev", Version: "v1", Resource: "tasks"}
)

// CRDManager manages Tekton CRD resources via the dynamic client.
type CRDManager struct {
	dynamicClient dynamic.Interface
}

// NewCRDManager creates a new CRDManager.
func NewCRDManager(dynamicClient dynamic.Interface) *CRDManager {
	return &CRDManager{dynamicClient: dynamicClient}
}

// CreatePipelineRun creates a PipelineRun from YAML bytes.
func (m *CRDManager) CreatePipelineRun(ctx context.Context, namespace string, yamlData []byte) (*unstructured.Unstructured, error) {
	return m.createResource(ctx, PipelineRunGVR, namespace, yamlData)
}

// CreateTaskRun creates a TaskRun from YAML bytes.
func (m *CRDManager) CreateTaskRun(ctx context.Context, namespace string, yamlData []byte) (*unstructured.Unstructured, error) {
	return m.createResource(ctx, TaskRunGVR, namespace, yamlData)
}

// GetPipelineRun gets a PipelineRun by name.
func (m *CRDManager) GetPipelineRun(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return m.dynamicClient.Resource(PipelineRunGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetTaskRun gets a TaskRun by name.
func (m *CRDManager) GetTaskRun(ctx context.Context, namespace, name string) (*unstructured.Unstructured, error) {
	return m.dynamicClient.Resource(TaskRunGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

// DeletePipelineRun deletes a PipelineRun.
func (m *CRDManager) DeletePipelineRun(ctx context.Context, namespace, name string) error {
	return m.dynamicClient.Resource(PipelineRunGVR).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// DeleteTaskRun deletes a TaskRun.
func (m *CRDManager) DeleteTaskRun(ctx context.Context, namespace, name string) error {
	return m.dynamicClient.Resource(TaskRunGVR).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// CancelPipelineRun patches a PipelineRun to cancel it.
func (m *CRDManager) CancelPipelineRun(ctx context.Context, namespace, name string) error {
	return m.cancelRun(ctx, PipelineRunGVR, namespace, name)
}

// CancelTaskRun patches a TaskRun to cancel it.
func (m *CRDManager) CancelTaskRun(ctx context.Context, namespace, name string) error {
	return m.cancelRun(ctx, TaskRunGVR, namespace, name)
}

// ListPipelineRuns lists PipelineRuns with a label selector.
func (m *CRDManager) ListPipelineRuns(ctx context.Context, namespace, labelSelector string) ([]unstructured.Unstructured, error) {
	list, err := m.dynamicClient.Resource(PipelineRunGVR).Namespace(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pipeline runs: %w", err)
	}
	return list.Items, nil
}

// ParseRunStatus extracts RunStatus from an unstructured PipelineRun/TaskRun.
func (m *CRDManager) ParseRunStatus(obj *unstructured.Unstructured) (*RunStatus, error) {
	rs := &RunStatus{
		Name: obj.GetName(),
	}

	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return nil, fmt.Errorf("failed to get conditions: %w", err)
	}
	if found && len(conditions) > 0 {
		cond, ok := conditions[0].(map[string]interface{})
		if ok {
			rs.Status = mapTektonStatus(
				getNestedString(cond, "type"),
				getNestedString(cond, "status"),
				getNestedString(cond, "reason"),
			)
			rs.Message = getNestedString(cond, "message")
		}
	} else {
		rs.Status = "pending"
	}

	if startStr, ok, _ := unstructured.NestedString(obj.Object, "status", "startTime"); ok {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			rs.StartedAt = &t
		}
	}
	if complStr, ok, _ := unstructured.NestedString(obj.Object, "status", "completionTime"); ok {
		if t, err := time.Parse(time.RFC3339, complStr); err == nil {
			rs.FinishedAt = &t
		}
	}

	steps, found, _ := unstructured.NestedSlice(obj.Object, "status", "steps")
	if found {
		for _, s := range steps {
			stepMap, ok := s.(map[string]interface{})
			if !ok {
				continue
			}
			rs.Steps = append(rs.Steps, StepStatus{
				Name:      getNestedString(stepMap, "name"),
				Container: getNestedString(stepMap, "container"),
			})
		}
	}

	return rs, nil
}

func (m *CRDManager) createResource(ctx context.Context, gvr schema.GroupVersionResource, namespace string, yamlData []byte) (*unstructured.Unstructured, error) {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(yamlData), 4096)
	var obj unstructured.Unstructured
	if err := decoder.Decode(&obj); err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("empty YAML input")
		}
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	result, err := m.dynamicClient.Resource(gvr).Namespace(namespace).Create(ctx, &obj, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create %s: %w", gvr.Resource, err)
	}
	return result, nil
}

func (m *CRDManager) cancelRun(ctx context.Context, gvr schema.GroupVersionResource, namespace, name string) error {
	patch := map[string]interface{}{
		"spec": map[string]interface{}{
			"status": "CancelledRunFinally",
		},
	}
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return fmt.Errorf("failed to marshal cancel patch: %w", err)
	}

	_, err = m.dynamicClient.Resource(gvr).Namespace(namespace).Patch(
		ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to cancel %s/%s: %w", gvr.Resource, name, err)
	}
	return nil
}

func mapTektonStatus(condType, condStatus, reason string) string {
	if condType != "Succeeded" {
		return "running"
	}
	switch condStatus {
	case "True":
		return "succeeded"
	case "False":
		if reason == "TaskRunCancelled" || reason == "PipelineRunCancelled" || reason == "Cancelled" {
			return "cancelled"
		}
		return "failed"
	default:
		return "running"
	}
}

func getNestedString(obj map[string]interface{}, key string) string {
	val, ok := obj[key]
	if !ok {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		return ""
	}
	return s
}
