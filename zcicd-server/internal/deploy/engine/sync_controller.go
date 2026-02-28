package engine

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

// SyncController handles sync operations on Argo CD Applications.
type SyncController struct {
	client        dynamic.Interface
	argoNamespace string
}

// NewSyncController creates a new SyncController.
func NewSyncController(dynamicClient dynamic.Interface, argoNamespace string) *SyncController {
	return &SyncController{
		client:        dynamicClient,
		argoNamespace: argoNamespace,
	}
}

// TriggerSync patches the Application with a sync operation annotation.
func (s *SyncController) TriggerSync(ctx context.Context, appName string, revision string) (*SyncResult, error) {
	patch := fmt.Sprintf(
		`{"metadata":{"annotations":{"argocd.argoproj.io/refresh":"hard"}},"operation":{"initiatedBy":{"username":"zcicd"},"sync":{"revision":"%s"}}}`,
		revision,
	)

	obj, err := s.client.Resource(argoAppGVR).Namespace(s.argoNamespace).
		Patch(ctx, appName, types.MergePatchType, []byte(patch), metav1.PatchOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to trigger sync for %s: %w", appName, err)
	}

	return extractSyncResult(obj), nil
}

// GetSyncStatus reads the current sync status from the Application.
func (s *SyncController) GetSyncStatus(ctx context.Context, appName string) (*SyncResult, error) {
	obj, err := s.client.Resource(argoAppGVR).Namespace(s.argoNamespace).
		Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get sync status for %s: %w", appName, err)
	}
	return extractSyncResult(obj), nil
}

// WaitForSync polls until the sync operation completes or the timeout is reached.
func (s *SyncController) WaitForSync(ctx context.Context, appName string, timeout time.Duration) (*SyncResult, error) {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return nil, fmt.Errorf("timeout waiting for sync of %s after %v", appName, timeout)
			}
			result, err := s.GetSyncStatus(ctx, appName)
			if err != nil {
				return nil, err
			}
			if result.Status == "Synced" && result.Health != "Progressing" {
				return result, nil
			}
		}
	}
}

// extractSyncResult parses sync and health status from an unstructured Application.
func extractSyncResult(obj *unstructured.Unstructured) *SyncResult {
	result := &SyncResult{}
	result.Status, _, _ = unstructured.NestedString(obj.Object, "status", "sync", "status")
	result.Health, _, _ = unstructured.NestedString(obj.Object, "status", "health", "status")
	result.Revision, _, _ = unstructured.NestedString(obj.Object, "status", "sync", "revision")
	result.Message, _, _ = unstructured.NestedString(obj.Object, "status", "operationState", "message")
	result.StartedAt, _, _ = unstructured.NestedString(obj.Object, "status", "operationState", "startedAt")
	result.FinishedAt, _, _ = unstructured.NestedString(obj.Object, "status", "operationState", "finishedAt")
	return result
}
