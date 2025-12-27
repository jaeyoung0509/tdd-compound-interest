CREATE TABLE outbox_messages (
    id             CHAR(26)    PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id   CHAR(26)    NOT NULL,
    event_type     VARCHAR(120) NOT NULL,
    payload        JSONB       NOT NULL,
    occurred_at    TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at   TIMESTAMPTZ NULL,
    attempts       INTEGER     NOT NULL DEFAULT 0
);

CREATE INDEX idx_outbox_unpublished ON outbox_messages (created_at) WHERE published_at IS NULL;
CREATE INDEX idx_outbox_aggregate ON outbox_messages (aggregate_type, aggregate_id);

COMMENT ON TABLE outbox_messages IS 'Outbox messages for domain event publication';
