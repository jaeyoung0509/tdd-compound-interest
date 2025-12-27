package payment

import "errors"

var (
	ErrInvalidUserID            = errors.New("invalid user id")
	ErrInvalidDueDate           = errors.New("invalid due date")
	ErrInvalidPaidAt            = errors.New("invalid paid at")
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentAlreadyPaid       = errors.New("payment already paid")
	ErrPaymentAlreadyOverdue    = errors.New("payment already overdue")
	ErrPaidPaymentCannotOverdue = errors.New("paid payment cannot be marked overdue")
	ErrInvalidOverdueArgs       = errors.New("invalid overdue args")
	ErrOverduePeriodTooLong     = errors.New("overdue period exceeds limit")
)
