package payment

import "errors"

var (
	ErrInvalidUserID            = errors.New("invalid user id")
	ErrInvalidAmount            = errors.New("invalid amount")
	ErrInvalidDueDate           = errors.New("invalid due date")
	ErrDueDateInPast            = errors.New("due date in past")
	ErrInvalidPaidAt            = errors.New("invalid paid at")
	ErrPaidBeforeDueDate        = errors.New("paid before due date")
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrPaymentAlreadyPaid       = errors.New("payment already paid")
	ErrPaymentAlreadyOverdue    = errors.New("payment already overdue")
	ErrPaidPaymentCannotOverdue = errors.New("paid payment cannot be marked overdue")
	ErrInvalidOverdueArgs       = errors.New("invalid overdue args")
	ErrOverduePeriodTooLong     = errors.New("overdue period exceeds limit")
)
