CREATE TABLE payment_overdues (
    id CHAR(26) PRIMARY KEY, 
    payment_id CHAR(26) NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    is_overdue BOOLEAN NOT NULL DEFAULT TRUE, 
    days_overdue INTEGER NOT NULL,
    penalty NUMERIC NOT NULL, 
    penalty_currency VARCHAR(3) NOT NULL,
    calculated_at DATE NOT NULL,
    created_at TIMESTAMPZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPZ NOT NULL DEFAULT NOW(),
)


REATE INDEX idx_payment_overdues_payment_id ON payment_overdues(payment_id);
CREATE INDEX idx_payment_overdues_calculated_at ON payment_overdues(calculated_at DESC);

COMMENT ON TABLE payment_overdues IS 'Immutable snapshots of overdue calculations (append-only history)';