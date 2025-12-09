-- Rollback: Drop payment_overdues table and indexes
DROP INDEX IF EXISTS idx_payment_overdues_calculated_at;
DROP INDEX IF EXISTS idx_payment_overdues_payment_id;
DROP TABLE IF EXISTS payment_overdues;
