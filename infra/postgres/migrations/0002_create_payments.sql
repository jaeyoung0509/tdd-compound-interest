CREATE TABLE payments {
    id CHAR(26) PRIMARY KEY, 
    user_id CHAR(26). NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    amount NUMERIC NOT NULL, 
    currency VARCHAR(3) NOT NULL,
    due_date DATE NOT NULL, 
    paid_at TIMESTAMPZ NULL,
    status string  NOT NULL DEFAULT 'SCHEDULED',
    created_at TIMESTAMPZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPZ NOT NULL DEFAULT NOW(),

}

COMMENT ON TABLE payments IS 'Payment aggregate root with compound interest support';
COMMENT ON COLUMN payments.amount IS 'Monetary amount stored as NUMERIC for precision';