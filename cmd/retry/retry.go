package retry

import (
	"context"
	"fmt"
	"time"
)

func Retry(ctx context.Context, attempts int, sleep time.Duration, f func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(sleep)
			sleep *= 2
		}

		err = f()

		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("after %d attempts, failed with %v", attempts, err)
}
