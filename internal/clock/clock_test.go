package clock_test

import (
	"testing"
	"time"

	"github.com/mercadolibre/fury_acq-visa-clearing/internal/platform/clock"
	"github.com/stretchr/testify/require"
)

func TestNowUTCTime(t *testing.T) {
	expectNow := time.Now().UTC()
	actualNow := clock.New().Now()

	require.EqualValues(t, expectNow.Format(time.RFC3339), actualNow.Format(time.RFC3339))
}
