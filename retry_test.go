package retry_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/vslpsl/retry"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	testError := errors.New("test-error")
	count := 0
	err := retry.Retry(
		retry.Execute(func() error {
			count++
			return testError
		}),
		retry.Inspect(func(r retry.Retrier, err error) {
			if r.PassedIterationCount() == 3 {
				r.Stop()
				return
			}
			r.SetDelay(time.Second)
		}),
	)

	require.ErrorIs(t, err, testError)
	require.Equal(t, 3, count)
}
