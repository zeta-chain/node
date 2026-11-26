package retry

import (
	"errors"
	"strings"
	"testing"

	"github.com/cenkalti/backoff/v4"
)

func TestDoTypedWithBackoffPreservesRetryLimit(t *testing.T) {
	want := errors.New("temporary failure")
	value, err := DoTypedWithBackoff(func() (int, error) {
		return 42, Retry(want)
	}, &backoff.StopBackOff{})
	var retried errRetryable
	if value != 42 || !errors.As(err, &retried) || !errors.Is(retried.error, want) || !strings.Contains(err.Error(), "retry limit exceeded") {
		t.Fatalf("value=%d err=%v; want last value and wrapped retry-limit error", value, err)
	}
}
