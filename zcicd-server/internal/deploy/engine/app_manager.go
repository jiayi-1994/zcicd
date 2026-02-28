package engine

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var argoAppGVR = schema.GroupVersionResource{
	Group:    "argoproj.io",
	Version:  "v1alpha1",
	Resource: "applications",
}

// AppManager manages Argo CD Application CRs via the dynamic client.
type AppManager struct {
	client        dynamic.Interface
	argoNamespace string
}

// NewAppManager creates a new AppManager.
func NewAppManager(dynamicClient dynamic.Interface, argoNamespace string) *AppManager {
	return &AppManager{
		client:        dynamicClient,
		argoNamespace: argoNamespace,
	}
}

// CreateApp creates an Argo CD Application CR.
func (m *AppManager) CreateApp(ctx context.Context, app ArgoApp) error {
	obj := buildApplicationCR(app, m.argoNamespace)
	_, err := m.client.Resource(argoAppGVR).Namespace(m.argoNamespace).
		Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create argo application %s: %w", app.Name, err)
	}
	return nil
}

// UpdateApp updates an existing Argo CD Application CR.
func (m *AppManager) UpdateApp(ctx context.Context, app ArgoApp) error {
	obj := buildApplicationCR(app, m.argoNamespace)

	existing, err := m.client.Resource(argoAppGVR).Namespace(m.argoNamespace).
		Get(ctx, app.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get argo application %s: %w", app.Name, err)
	}
	obj.SetResourceVersion(existing.GetResourceVersion())

	_, err = m.client.Resource(argoAppGVR).Namespace(m.argoNamespace).
		Update(ctx, obj, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update argo application %s: %w", app.Name, err)
	}
	return nil
}

// DeleteApp deletes an Argo CD Application CR.
func (m *AppManager) DeleteApp(ctx context.Context, name string) error {
	err := m.client.Resource(argoAppGVR).Namespace(m.argoNamespace).
		Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete argo application %s: %w", name, err)
	}
	return nil
}

// GetApp retrieves an Argo CD Application and parses its status.
func (m *AppManager) GetApp(ctx context.Context, name string) (*AppStatus, error) {
	obj, err := m.client.Resource(argoAppGVR).Namespace(m.argoNamespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get argo application %s: %w", name, err)
	}
	return parseAppStatus(obj), nil
}

// GetResourceTree returns the managed resources from the Application status.
func (m *AppManager) GetResourceTree(ctx context.Context, name string) (*ResourceTree, error) {
	obj, err := m.client.Resource(argoAppGVR).Namespace(m.argoNamespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get argo application %s: %w", name, err)
	}

	resources, _, _ := unstructured.NestedSlice(obj.Object, "status", "resources")
	tree := &ResourceTree{}
	for _, r := range resources {
		rm, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		tree.Nodes = append(tree.Nodes, parseResourceNode(rm))
	}
	return tree, nil
}

// buildApplicationCR constructs an unstructured Argo CD Application CR.
func buildApplicationCR(app ArgoApp, argoNamespace string) *unstructured.Unstructured {
	project := app.Project
	if project == "" {
		project = "default"
	}
	destServer := app.DestServer
	if destServer == "" {
		destServer = "https://kubernetes.default.svc"
	}
	ns := app.Namespace
	if ns == "" {
		ns = argoNamespace
	}

	source := map[string]interface{}{
		"repoURL":        app.RepoURL,
		"targetRevision": app.TargetRevision,
		"path":           app.Path,
	}
	if len(app.ValuesOverride) > 0 {
		source["helm"] = map[string]interface{}{
			"values": mapToYAMLString(app.ValuesOverride),
		}
	}

	spec := map[string]interface{}{
		"project": project,
		"source":  source,
		"destination": map[string]interface{}{
			"server":    destServer,
			"namespace": app.DestNamespace,
		},
	}

	if app.AutoSync {
		automated := map[string]interface{}{
			"selfHeal": app.SelfHeal,
			"prune":    app.Prune,
		}
		spec["syncPolicy"] = map[string]interface{}{
			"automated": automated,
		}
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1",
			"kind":       "Application",
			"metadata": map[string]interface{}{
				"name":      app.Name,
				"namespace": ns,
			},
			"spec": spec,
		},
	}
	return obj
}

// parseAppStatus extracts AppStatus from an unstructured Application.
func parseAppStatus(obj *unstructured.Unstructured) *AppStatus {
	status := &AppStatus{}
	status.SyncStatus, _, _ = unstructured.NestedString(obj.Object, "status", "sync", "status")
	status.HealthStatus, _, _ = unstructured.NestedString(obj.Object, "status", "health", "status")
	status.Revision, _, _ = unstructured.NestedString(obj.Object, "status", "sync", "revision")
	status.Message, _, _ = unstructured.NestedString(obj.Object, "status", "health", "message")

	resources, _, _ := unstructured.NestedSlice(obj.Object, "status", "resources")
	for _, r := range resources {
		rm, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		status.Resources = append(status.Resources, parseResourceNode(rm))
	}
	return status
}

// parseResourceNode extracts a ResourceNode from a map.
func parseResourceNode(rm map[string]interface{}) ResourceNode {
	getString := func(key string) string {
		v, _ := rm[key].(string)
		return v
	}
	health := ""
	if h, ok := rm["health"].(map[string]interface{}); ok {
		health, _ = h["status"].(string)
	}
	return ResourceNode{
		Group:     getString("group"),
		Kind:      getString("kind"),
		Namespace: getString("namespace"),
		Name:      getString("name"),
		Status:    getString("status"),
		Health:    health,
	}
}

// mapToYAMLString converts a map to a simple YAML-like key: value string.
func mapToYAMLString(m map[string]interface{}) string {
	result := ""
	for k, v := range m {
		result += fmt.Sprintf("%s: %v\n", k, v)
	}
	return result
}
