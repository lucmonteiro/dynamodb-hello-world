package clock_test

import (
	"dynamo-hello-world/internal/clock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNowUTCTime(t *testing.T) {
	expectNow := time.Now().UTC()
	actualNow := clock.New().Now()

	require.EqualValues(t, expectNow.Format(time.RFC3339), actualNow.Format(time.RFC3339))
}
