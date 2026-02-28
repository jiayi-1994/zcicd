package engine

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var rolloutGVR = schema.GroupVersionResource{
	Group:    "argoproj.io",
	Version:  "v1alpha1",
	Resource: "rollouts",
}

type RolloutController struct {
	client dynamic.Interface
	ns     string
}

func NewRolloutController(client dynamic.Interface, ns string) *RolloutController {
	return &RolloutController{client: client, ns: ns}
}

func (r *RolloutController) GetStatus(ctx context.Context, name string) (*RolloutStatus, error) {
	obj, err := r.client.Resource(rolloutGVR).Namespace(r.ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get rollout %s: %w", name, err)
	}
	return parseRolloutStatus(obj)
}

func (r *RolloutController) Promote(ctx context.Context, name string) error {
	obj, err := r.client.Resource(rolloutGVR).Namespace(r.ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get rollout %s: %w", name, err)
	}
	// Set status.promoteFull to trigger promotion
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["rollout.argoproj.io/promote"] = "full"
	obj.SetAnnotations(annotations)
	_, err = r.client.Resource(rolloutGVR).Namespace(r.ns).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func (r *RolloutController) Abort(ctx context.Context, name string) error {
	obj, err := r.client.Resource(rolloutGVR).Namespace(r.ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get rollout %s: %w", name, err)
	}
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["rollout.argoproj.io/abort"] = "true"
	obj.SetAnnotations(annotations)
	_, err = r.client.Resource(rolloutGVR).Namespace(r.ns).Update(ctx, obj, metav1.UpdateOptions{})
	return err
}

func parseRolloutStatus(obj *unstructured.Unstructured) (*RolloutStatus, error) {
	statusRaw, ok := obj.Object["status"]
	if !ok {
		return &RolloutStatus{Phase: "Unknown"}, nil
	}
	data, err := json.Marshal(statusRaw)
	if err != nil {
		return nil, err
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	s := &RolloutStatus{}
	if v, ok := raw["phase"].(string); ok {
		s.Phase = v
	}
	if v, ok := raw["message"].(string); ok {
		s.Message = v
	}
	if v, ok := raw["stableRS"].(string); ok {
		s.StableRevision = v
	}
	if v, ok := raw["currentStepIndex"].(float64); ok {
		s.CurrentStep = int(v)
	}
	return s, nil
}
