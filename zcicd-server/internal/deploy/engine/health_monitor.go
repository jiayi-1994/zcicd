package engine

import (
	"context"
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

// HealthMonitor watches Argo CD Applications for health status changes.
type HealthMonitor struct {
	client        dynamic.Interface
	argoNamespace string
	stopCh        chan struct{}
	mu            sync.Mutex
	running       bool
}

// NewHealthMonitor creates a new HealthMonitor.
func NewHealthMonitor(dynamicClient dynamic.Interface, argoNamespace string) *HealthMonitor {
	return &HealthMonitor{
		client:        dynamicClient,
		argoNamespace: argoNamespace,
	}
}

// Start begins watching Argo CD Applications and calls the callback on status changes.
func (h *HealthMonitor) Start(ctx context.Context, callback func(appName string, status AppStatus)) {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return
	}
	h.stopCh = make(chan struct{})
	h.running = true
	h.mu.Unlock()

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		h.client, 0, h.argoNamespace, nil,
	)

	informer := factory.ForResource(argoAppGVR).Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			u, ok := newObj.(*unstructured.Unstructured)
			if !ok {
				return
			}
			appName := u.GetName()
			status := parseAppStatus(u)
			callback(appName, *status)
		},
	})

	go informer.Run(h.stopCh)

	// Wait for context cancellation to stop
	go func() {
		select {
		case <-ctx.Done():
			h.Stop()
		case <-h.stopCh:
		}
	}()
}

// Stop stops the informer.
func (h *HealthMonitor) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.running {
		close(h.stopCh)
		h.running = false
	}
}

// GetHealth performs a one-shot health check for an application.
func (h *HealthMonitor) GetHealth(ctx context.Context, appName string) (string, error) {
	obj, err := h.client.Resource(argoAppGVR).Namespace(h.argoNamespace).
		Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get health for %s: %w", appName, err)
	}
	health, _, _ := unstructured.NestedString(obj.Object, "status", "health", "status")
	if health == "" {
		health = "Unknown"
	}
	return health, nil
}
