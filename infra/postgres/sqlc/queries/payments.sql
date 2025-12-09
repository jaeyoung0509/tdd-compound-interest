-- name: GetPayment :one 
SELECT id, user_id, amount, currency, due_date, paid_at, status, created_at, updated_at 
FROM payments 
WHERE id = $1; 


-- name: UpsertPayment :exec
INSERT INTO payments (id, user_id, amount, currency, due_date, paid_at, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (id) DO UPDATE SET 
    amount = EXCLUDED.amount,
    paid_at = EXCLUDED.paid_at,
    status = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at; 


-- name: GetLatestOverdue :one
SELECT id, payment_id, is_overdue, days_overdue, penalty, penalty_currency, calculated_at, created_at
FROM payment_overdues
WHERE payment_id = $1
ORDER BY calculated_at DESC
LIMIT 1;

-- name: InsertOverdue :exec
INSERT INTO payment_overdues (id, payment_id, is_overdue, days_overdue, penalty, penalty_currency, calculated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);