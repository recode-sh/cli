package stepper

import (
	"fmt"

	"github.com/briandowns/spinner"
)

type Step struct {
	spin            *spinner.Spinner
	removeAfterDone bool
}

func (s *Step) Done() {
	s.spin.Stop()

	if !s.removeAfterDone {
		fmt.Println(s.spin.Prefix + "... done")
	}
}
