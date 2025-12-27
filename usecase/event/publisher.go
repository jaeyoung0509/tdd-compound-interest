package event

import (
	"context"

	"github.com/jaeyoung0509/compound-interest/domain/shared"
)

// Publisher persists or dispatches domain events (outbox-backed in production).
type Publisher interface {
	Publish(ctx context.Context, events ...shared.DomainEvent) error
}

// NoopPublisher ignores events, useful for tests.
type NoopPublisher struct{}

func (NoopPublisher) Publish(ctx context.Context, events ...shared.DomainEvent) error {
	return nil
}
