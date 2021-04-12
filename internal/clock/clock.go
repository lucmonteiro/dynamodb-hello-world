package clock

import "time"

type Clock interface {
	Now() time.Time
}

type clock struct{}

func New() Clock {
	return &clock{}
}

func (clock) Now() time.Time {
	return time.Now().UTC()
}
