package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClient wraps the typed and dynamic Kubernetes clients.
type K8sClient struct {
	Clientset     kubernetes.Interface
	DynamicClient dynamic.Interface
	Config        *rest.Config
}

// NewK8sClient creates a K8sClient from a kubeconfig path.
// If kubeconfig is empty, it falls back to in-cluster configuration.
func NewK8sClient(kubeconfig string) (*K8sClient, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s config: %w", err)
	}

	return NewK8sClientFromConfig(config)
}

// NewK8sClientFromConfig creates a K8sClient from an existing rest.Config.
func NewK8sClientFromConfig(config *rest.Config) (*K8sClient, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &K8sClient{
		Clientset:     clientset,
		DynamicClient: dynClient,
		Config:        config,
	}, nil
}

// GetNamespaces returns a list of all namespace names in the cluster.
func (c *K8sClient) GetNamespaces(ctx context.Context) ([]string, error) {
	nsList, err := c.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	names := make([]string, 0, len(nsList.Items))
	for _, ns := range nsList.Items {
		names = append(names, ns.Name)
	}
	return names, nil
}

// GetPods returns all pods in the given namespace.
func (c *K8sClient) GetPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	podList, err := c.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}
	return podList.Items, nil
}

// GetDeployments returns all deployments in the given namespace.
func (c *K8sClient) GetDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	depList, err := c.Clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}
	return depList.Items, nil
}

// ApplyYAML applies a YAML manifest using the dynamic client.
func (c *K8sClient) ApplyYAML(ctx context.Context, namespace string, yaml []byte) error {
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(yaml), 4096)

	for {
		var obj unstructured.Unstructured
		if err := decoder.Decode(&obj); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode YAML: %w", err)
		}

		gvk := obj.GroupVersionKind()
		gvr := schema.GroupVersionResource{
			Group:    gvk.Group,
			Version:  gvk.Version,
			Resource: pluralizeKind(gvk.Kind),
		}

		if namespace != "" {
			obj.SetNamespace(namespace)
		}

		ns := obj.GetNamespace()
		var err error
		if ns != "" {
			_, err = c.DynamicClient.Resource(gvr).Namespace(ns).
				Apply(ctx, obj.GetName(), &obj, metav1.ApplyOptions{FieldManager: "zcicd"})
		} else {
			_, err = c.DynamicClient.Resource(gvr).
				Apply(ctx, obj.GetName(), &obj, metav1.ApplyOptions{FieldManager: "zcicd"})
		}
		if err != nil {
			return fmt.Errorf("failed to apply %s/%s: %w", gvk.Kind, obj.GetName(), err)
		}
	}

	return nil
}

// DeleteResource deletes a single resource by apiVersion, kind, and name.
func (c *K8sClient) DeleteResource(ctx context.Context, namespace, apiVersion, kind, name string) error {
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return fmt.Errorf("failed to parse apiVersion %s: %w", apiVersion, err)
	}

	gvr := schema.GroupVersionResource{
		Group:    gv.Group,
		Version:  gv.Version,
		Resource: pluralizeKind(kind),
	}

	if namespace != "" {
		return c.DynamicClient.Resource(gvr).Namespace(namespace).
			Delete(ctx, name, metav1.DeleteOptions{})
	}
	return c.DynamicClient.Resource(gvr).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetLogs returns the logs for a specific pod container.
func (c *K8sClient) GetLogs(ctx context.Context, namespace, podName, container string, tailLines int64) (string, error) {
	opts := &corev1.PodLogOptions{
		TailLines: &tailLines,
	}
	if container != "" {
		opts.Container = container
	}

	req := c.Clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get logs for pod %s/%s: %w", namespace, podName, err)
	}
	defer stream.Close()

	buf, err := io.ReadAll(stream)
	if err != nil {
		return "", fmt.Errorf("failed to read log stream: %w", err)
	}
	return string(buf), nil
}

// pluralizeKind converts a Kind (e.g. "Deployment") to its plural resource name.
func pluralizeKind(kind string) string {
	lower := strings.ToLower(kind)
	switch {
	case strings.HasSuffix(lower, "s"):
		return lower + "es"
	case strings.HasSuffix(lower, "y"):
		return lower[:len(lower)-1] + "ies"
	default:
		return lower + "s"
	}
}
