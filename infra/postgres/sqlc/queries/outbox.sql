-- name: InsertOutbox :exec
INSERT INTO outbox_messages (
    id, aggregate_type, aggregate_id, event_type, payload, occurred_at
)
VALUES ($1, $2, $3, $4, $5, $6);
