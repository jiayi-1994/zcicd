package engine

import (
	"bufio"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// LogCollector streams and collects logs from Tekton run pods.
type LogCollector struct {
	k8sClient kubernetes.Interface
	rdb       *redis.Client
}

// NewLogCollector creates a new LogCollector.
func NewLogCollector(k8sClient kubernetes.Interface, rdb *redis.Client) *LogCollector {
	return &LogCollector{
		k8sClient: k8sClient,
		rdb:       rdb,
	}
}

// StreamLogs streams logs from a pod container to a Redis Pub/Sub channel.
// Channel format: "build:logs:{runID}"
func (c *LogCollector) StreamLogs(ctx context.Context, namespace, podName, container, runID string) error {
	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    true,
	}

	stream, err := c.k8sClient.CoreV1().Pods(namespace).GetLogs(podName, opts).Stream(ctx)
	if err != nil {
		return fmt.Errorf("failed to open log stream for pod %s/%s: %w", namespace, podName, err)
	}
	defer stream.Close()

	channel := fmt.Sprintf("build:logs:%s", runID)
	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line := scanner.Text()
			if err := c.rdb.Publish(ctx, channel, line).Err(); err != nil {
				return fmt.Errorf("failed to publish log line: %w", err)
			}
		}
	}

	return scanner.Err()
}

// GetLogs gets completed logs from a pod (non-streaming).
func (c *LogCollector) GetLogs(ctx context.Context, namespace, podName, container string, tailLines int64) (string, error) {
	opts := &corev1.PodLogOptions{
		Container: container,
		TailLines: &tailLines,
	}

	result := c.k8sClient.CoreV1().Pods(namespace).GetLogs(podName, opts).Do(ctx)
	raw, err := result.Raw()
	if err != nil {
		return "", fmt.Errorf("failed to get logs for pod %s/%s: %w", namespace, podName, err)
	}
	return string(raw), nil
}

// SubscribeLogs subscribes to a Redis Pub/Sub channel for log streaming.
func (c *LogCollector) SubscribeLogs(ctx context.Context, runID string) (<-chan string, error) {
	channel := fmt.Sprintf("build:logs:%s", runID)
	sub := c.rdb.Subscribe(ctx, channel)

	// Wait for subscription confirmation
	if _, err := sub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", channel, err)
	}

	out := make(chan string, 100)
	go func() {
		defer close(out)
		defer sub.Close()

		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- msg.Payload:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out, nil
}

// ArchiveLogs saves completed logs to object storage.
// TODO: integrate with MinIO storage package when available.
func (c *LogCollector) ArchiveLogs(ctx context.Context, runID string, logs string, storagePath string) error {
	// Placeholder: will write logs to MinIO via the storage package.
	// For now, store in Redis with a TTL as a temporary measure.
	key := fmt.Sprintf("build:logs:archive:%s", runID)
	if err := c.rdb.Set(ctx, key, logs, 0).Err(); err != nil {
		return fmt.Errorf("failed to archive logs for run %s: %w", runID, err)
	}
	return nil
}
