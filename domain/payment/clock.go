package payment

import "time"

// Clock abstracts current time retrieval for testability.
type Clock interface {
	Now() time.Time
}

// RealClock uses time.Now and is intended for production wiring.
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

// FixedClock returns a preset time, useful for tests.
type FixedClock struct {
	NowTime time.Time
}

func (c FixedClock) Now() time.Time {
	return c.NowTime
}

// DailyRateProvider supplies a daily basis points rate at a given time.
type DailyRateProvider interface {
	DailyRateBPS(at time.Time) (int64, error)
}

// StaticDailyRate provides a constant rate regardless of time, handy for tests.
type StaticDailyRate struct {
	BPS int64
	Err error
}

func (r StaticDailyRate) DailyRateBPS(time.Time) (int64, error) {
	return r.BPS, r.Err
}
