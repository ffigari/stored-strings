package clock

import (
	"time"
)

type clock struct{}

func New() *clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}
