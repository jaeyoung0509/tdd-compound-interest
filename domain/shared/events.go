package shared

import "time"

// DomainEvent captures metadata required for outbox publication.
type DomainEvent interface {
	EventType() string
	AggregateType() string
	AggregateID() string
	OccurredAt() time.Time
}
