package interfaces

import "time"

type Sleeper interface {
	Sleep(d time.Duration)
}
