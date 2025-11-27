package payment

type Status string

const (
	StatusScheduled Status = "SCHEDULED"
	StatusPaid      Status = "PAID"
	StatusOverdue   Status = "OVERDUE"
)
