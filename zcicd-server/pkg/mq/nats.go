package mq

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/zcicd/zcicd-server/pkg/config"
)

type Client struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewNATSClient(cfg *config.Config) (*Client, error) {
	nc, err := nats.Connect(cfg.NATS.URL,
		nats.MaxReconnects(cfg.NATS.MaxReconnects),
		nats.ReconnectWait(nats.DefaultReconnectWait),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Ensure stream exists
	_, err = js.AddStream(&nats.StreamConfig{
		Name:      cfg.NATS.StreamName,
		Subjects:  []string{"zcicd.>"},
		Retention: nats.WorkQueuePolicy,
		MaxAge:    7 * 24 * 3600 * 1e9, // 7 days
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return &Client{conn: nc, js: js}, nil
}

func (c *Client) Publish(subject string, data []byte) error {
	_, err := c.js.Publish(subject, data)
	return err
}

func (c *Client) Subscribe(subject, consumer string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return c.js.Subscribe(subject, handler,
		nats.Durable(consumer),
		nats.AckWait(30*1e9),
		nats.MaxDeliver(5),
	)
}

func (c *Client) Close() {
	c.conn.Close()
}
