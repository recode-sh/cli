package system

import (
	"time"
)

type Sleeper struct{}

func NewSleeper() Sleeper {
	return Sleeper{}
}

func (Sleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}
