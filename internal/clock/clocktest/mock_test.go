package clocktest_test

import (
	"testing"
	"time"

	"github.com/mercadolibre/fury_acq-visa-clearing/internal/platform/clock/clocktest"
	"github.com/stretchr/testify/require"
)

func TestClock_SpecificMockedNow(t *testing.T) {
	nowExample := time.Date(2009, 10, 1, 11, 4, 5, 1, time.UTC)
	mock := clocktest.NewMockWithTime(nowExample)
	require.EqualValues(t, nowExample, mock.Now())

	mock.Add(10 * time.Second)
	require.EqualValues(t, nowExample.Add(10*time.Second), mock.Now())
}
