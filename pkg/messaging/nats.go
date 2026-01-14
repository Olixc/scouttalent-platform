package messaging

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// NATSConfig holds NATS configuration
type NATSConfig struct {
	URL string
}

// NewNATSClient creates a new NATS client connection
func NewNATSClient(cfg NATSConfig) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return nc, nil
}