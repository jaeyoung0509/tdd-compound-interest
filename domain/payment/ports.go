package payment

import "time"

// Clock abstracts time retrieval for deterministic tests.
type Clock interface {
	Now() time.Time
}

// DailyRateProvider supplies daily interest rate in basis points for a given date.
type DailyRateProvider interface {
	DailyRateBPS(at time.Time) (int64, error)
}

// FixedClock returns a fixed time, useful for tests.
type FixedClock struct {
	NowTime time.Time
}

func (c FixedClock) Now() time.Time {
	return c.NowTime
}

// StaticDailyRate always returns the configured BPS, useful for fixed-rate scenarios.
type StaticDailyRate struct {
	BPS int64
}

func (r StaticDailyRate) DailyRateBPS(time.Time) (int64, error) {
	return r.BPS, nil
}
