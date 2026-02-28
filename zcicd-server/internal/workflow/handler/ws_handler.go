package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	redisClient *redis.Client
}

func NewWSHandler(redisClient *redis.Client) *WSHandler {
	return &WSHandler{redisClient: redisClient}
}

func (h *WSHandler) HandleBuildLogs(c *gin.Context) {
	runID := c.Param("run_id")
	if runID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "run_id is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// Monitor client disconnect
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	// Subscribe to Redis Pub/Sub channel
	channel := fmt.Sprintf("build:logs:%s", runID)
	pubsub := h.redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	// Set a write deadline helper
	writeTimeout := 10 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Printf("websocket write failed for run %s: %v", runID, err)
				return
			}
		}
	}
}
