package cfg

import (
	"time"
)

type RegularInterval struct {
	// Fixed interval between each retry
	Interval time.Duration
	// Maximum number of attempts. When exceeded the retries stop.
	MaxAttempts uint
}

type RegularIntervalWithJitter struct {
	Interval    time.Duration
	Jitter      float64
	MaxAttempts uint
}

type ExponentialBackOff struct {
	// Backoff interval for the first retry.
	// If not set or set to 0, a default interval of 1s will be used.
	BaseInterval time.Duration
	// Maximum backoff interval between retries. Exponential backoff leads to interval increase.
	// This value is the cap of the interval. By default, there is no limit on the max interval.
	MaxInterval time.Duration
	// Coefficient used to calculate the next retry backoff interval.
	// The next retry interval is previous interval multiplied by this coefficient.
	// Must be larger than 1. Default is 2.0.
	// Use RegularInterval for cases where BackoffCoefficient is 1.
	BackoffCoefficient float64
	// Maximum number of attempts. When exceeded the retries stop.
	MaxAttempts uint
}

type ExponentialBackOffWithJitter struct {
	BaseInterval time.Duration
	MaxInterval  time.Duration
	Jitter       float64
	MaxAttempts  uint
}

type RandomizedInterval struct {
	MinInterval time.Duration
	MaxInterval time.Duration
	MaxAttempts uint
}

// Hybrid strategy follows a mix of two retry strategy
// If attempt is <= CutOff then RetryStrategy1 is followed
// If attempt is > CutOff then RetryStrategy2 is followed
type Hybrid struct {
	RetryStrategy1 *RetryParams
	RetryStrategy2 *RetryParams

	// Maximum number of retry attempts that can be made
	MaxAttempts uint

	// Cut off threshold to follow RetryStrategy2
	// Must be < MaxAttempts, for RetryStrategy2 to kick in
	CutOff uint
}
