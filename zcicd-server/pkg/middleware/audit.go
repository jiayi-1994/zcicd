package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zcicd/zcicd-server/pkg/mq"
)

// MQPublisher abstracts message queue publishing for testability.
type MQPublisher interface {
	Publish(subject string, data []byte) error
}

// AuditEntry is the payload recorded for each audited request.
type AuditEntry struct {
	UserID    string  `json:"user_id"`
	Username  string  `json:"username"`
	Method    string  `json:"method"`
	Path      string  `json:"path"`
	Status    int     `json:"status"`
	LatencyMs float64 `json:"latency_ms"`
	ClientIP  string  `json:"client_ip"`
	RequestID string  `json:"request_id"`
}

// AuditLog publishes an audit event after every request completes.
func AuditLog(mqClient MQPublisher) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		entry := AuditEntry{
			UserID:    c.GetString("user_id"),
			Username:  c.GetString("username"),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			LatencyMs: float64(time.Since(start).Milliseconds()),
			ClientIP:  c.ClientIP(),
			RequestID: c.GetString("request_id"),
		}

		event := mq.Event{
			EventID:     uuid.New().String(),
			EventType:   "audit.request",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			TriggeredBy: entry.Username,
			Payload:     entry,
		}

		go func() {
			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("[audit] marshal error: %v", err)
				return
			}
			if err := mqClient.Publish(mq.SubjectAuditLog, data); err != nil {
				log.Printf("[audit] publish error: %v", err)
			}
		}()
	}
}
