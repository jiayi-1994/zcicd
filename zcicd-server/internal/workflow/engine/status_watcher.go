package engine

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

// StatusCallback is called when a run status changes.
type StatusCallback func(runName string, status *RunStatus)

// StatusWatcher watches PipelineRun and TaskRun resources for status changes.
type StatusWatcher struct {
	dynamicClient dynamic.Interface
	namespace     string
	crdManager    *CRDManager
	mu            sync.RWMutex
	callbacks     map[string]StatusCallback
	stopCh        chan struct{}
}

// NewStatusWatcher creates a new StatusWatcher.
func NewStatusWatcher(dynamicClient dynamic.Interface, namespace string) *StatusWatcher {
	return &StatusWatcher{
		dynamicClient: dynamicClient,
		namespace:     namespace,
		crdManager:    NewCRDManager(dynamicClient),
		callbacks:     make(map[string]StatusCallback),
		stopCh:        make(chan struct{}),
	}
}

// Start begins watching PipelineRun and TaskRun resources.
func (w *StatusWatcher) Start(ctx context.Context) error {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		w.dynamicClient, 0, w.namespace, nil,
	)

	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			w.handleEvent(obj)
		},
		UpdateFunc: func(_, newObj interface{}) {
			w.handleEvent(newObj)
		},
	}

	pipelineRunInformer := factory.ForResource(PipelineRunGVR).Informer()
	taskRunInformer := factory.ForResource(TaskRunGVR).Informer()

	pipelineRunInformer.AddEventHandler(handler)
	taskRunInformer.AddEventHandler(handler)

	factory.Start(w.stopCh)
	factory.WaitForCacheSync(w.stopCh)

	go func() {
		<-ctx.Done()
		w.Stop()
	}()

	return nil
}

func (w *StatusWatcher) handleEvent(obj interface{}) {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return
	}

	runName := u.GetName()

	w.mu.RLock()
	cb, exists := w.callbacks[runName]
	w.mu.RUnlock()

	if !exists {
		return
	}

	status, err := w.crdManager.ParseRunStatus(u)
	if err != nil {
		return
	}

	cb(runName, status)
}

// RegisterCallback registers a callback for a specific run.
func (w *StatusWatcher) RegisterCallback(runName string, callback StatusCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.callbacks[runName] = callback
}

// UnregisterCallback removes a callback.
func (w *StatusWatcher) UnregisterCallback(runName string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.callbacks, runName)
}

// Stop stops the watcher.
func (w *StatusWatcher) Stop() {
	select {
	case <-w.stopCh:
		// already closed
	default:
		close(w.stopCh)
	}
}
