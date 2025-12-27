package repositories

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
	"github.com/jaeyoung0509/compound-interest/infra/postgres/sqlc/generated"
	"github.com/jaeyoung0509/compound-interest/usecase/event"
)

type OutboxPublisher struct {
	queries *generated.Queries
}

func NewOutboxPublisher(db generated.DBTX) *OutboxPublisher {
	return &OutboxPublisher{queries: generated.New(db)}
}

func (p *OutboxPublisher) Publish(ctx context.Context, events ...shared.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	for _, evt := range events {
		payloadBytes, err := json.Marshal(evt)
		if err != nil {
			return err
		}

		if err := p.queries.InsertOutbox(ctx, generated.InsertOutboxParams{
			ID:            shared.NewID().String(),
			AggregateType: evt.AggregateType(),
			AggregateID:   evt.AggregateID(),
			EventType:     evt.EventType(),
			Payload:       payloadBytes,
			OccurredAt:    pgtype.Timestamptz{Time: evt.OccurredAt(), Valid: true},
		}); err != nil {
			return err
		}
	}

	return nil
}

var _ event.Publisher = (*OutboxPublisher)(nil)
