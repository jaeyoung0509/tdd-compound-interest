package payment

import "github.com/jaeyoung0509/compound-interest/domain/shared"

var (
	ErrInvalidUserID            = shared.NewDomainError(shared.CodeValidation, "invalid user id")
	ErrInvalidAmount            = shared.NewDomainError(shared.CodeValidation, "invalid amount")
	ErrInvalidDueDate           = shared.NewDomainError(shared.CodeValidation, "invalid due date")
	ErrDueDateInPast            = shared.NewDomainError(shared.CodeValidation, "due date in past")
	ErrInvalidPaidAt            = shared.NewDomainError(shared.CodeValidation, "invalid paid at")
	ErrPaidBeforeDueDate        = shared.NewDomainError(shared.CodeConflict, "paid before due date")
	ErrPaymentNotFound          = shared.NewDomainError(shared.CodeNotFound, "payment not found")
	ErrPaymentAlreadyPaid       = shared.NewDomainError(shared.CodeConflict, "payment already paid")
	ErrPaymentAlreadyOverdue    = shared.NewDomainError(shared.CodeConflict, "payment already overdue")
	ErrPaidPaymentCannotOverdue = shared.NewDomainError(shared.CodeConflict, "paid payment cannot be marked overdue")
	ErrInvalidOverdueArgs       = shared.NewDomainError(shared.CodeValidation, "invalid overdue args")
	ErrOverduePeriodTooLong     = shared.NewDomainError(shared.CodeConflict, "overdue period exceeds limit")
)
