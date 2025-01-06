package retry

import (
	"errors"
	"time"
)

var ErrRetry = errors.New("retry")

func DoWithRetry[R any](unitOfWork func() (R, error), strategy RetryStrategy) (R, error) {
	for {
		r, err := unitOfWork()
		if errors.Is(err, ErrRetry) {
			sleep, sleepErr := strategy.NextBackoff()
			if sleepErr != nil {
				return *new(R), sleepErr
			}
			time.Sleep(sleep)
			continue
		}
		return r, err
	}
}
