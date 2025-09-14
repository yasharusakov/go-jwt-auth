package broker

import (
	"fmt"
	"log"
	"time"
	"user-service/internal/config"

	"github.com/nats-io/nats.go"
)

func NewNATS(cfg *config.NATSConfig, serviceName string) (*nats.Conn, error) {
	natsUrl := fmt.Sprintf("nats://%s:%s@%s:%s",
		cfg.NatsUser,
		cfg.NatsPassword,
		cfg.NatsHost,
		cfg.NatsPort,
	)

	nc, err := nats.Connect(
		natsUrl,
		nats.Name(serviceName),
		nats.Timeout(10*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Printf("Connected to NATS: %s", serviceName)

	return nc, nil
}
