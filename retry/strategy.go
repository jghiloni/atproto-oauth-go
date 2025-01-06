package retry

import (
	"errors"
	"math"
	"time"
)

var ErrTooManyRetries = errors.New("too many retries")

type RetryStrategy interface {
	NextBackoff() (time.Duration, error)
}

type LinearRetryStrategy struct {
	currentRetry int
	MaxRetries   int
	Backoff      time.Duration
	MaxBackoff   time.Duration
}

func (l *LinearRetryStrategy) NextBackoff() (time.Duration, error) {
	if l.currentRetry < 0 {
		l.currentRetry = 0
	}

	if l.currentRetry == math.MaxInt || l.currentRetry == l.MaxRetries {
		return 0, ErrTooManyRetries
	}

	nextBackoff := time.Duration(l.currentRetry+1) * l.Backoff

	if l.MaxBackoff > 0 {
		nextBackoff = min(nextBackoff, l.MaxBackoff)
	}

	l.currentRetry = l.currentRetry + 1
	return nextBackoff, nil
}

var _ RetryStrategy = &LinearRetryStrategy{}
