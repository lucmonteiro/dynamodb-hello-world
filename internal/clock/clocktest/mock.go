package clocktest

import "time"

type ClockMock struct {
	tim time.Time
}

func NewMockWithTime(t time.Time) *ClockMock {
	return &ClockMock{
		tim: t,
	}
}

func (ck *ClockMock) Now() time.Time {
	return ck.tim
}

func (ck *ClockMock) Add(duration time.Duration) {
	ck.tim = ck.tim.Add(duration)
}
