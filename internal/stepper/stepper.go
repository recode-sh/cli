package stepper

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/recode/stepper"
)

var currentStep *Step

type Stepper struct{}

func NewStepper() Stepper {
	return Stepper{}
}

func (s Stepper) startStep(
	step string,
	removeAfterDone bool,
	noNewLineAtStart bool,
) stepper.Step {

	if currentStep == nil && !noNewLineAtStart {
		fmt.Println("")
	}

	if currentStep != nil {
		currentStep.Done()
		currentStep = nil
	}

	bold := constants.Bold

	spin := spinner.New(spinner.CharSets[26], 400*time.Millisecond)
	spin.Prefix = bold(step)
	spin.Start()

	currentStep = &Step{
		spin:            spin,
		removeAfterDone: removeAfterDone,
	}

	return currentStep
}

func (s Stepper) StartStep(
	step string,
) stepper.Step {

	removeAfterDone := false
	noNewLineAtStart := false

	return s.startStep(
		step,
		removeAfterDone,
		noNewLineAtStart,
	)
}

func (s Stepper) StartTemporaryStep(
	step string,
) stepper.Step {

	removeAfterDone := true
	noNewLineAtStart := false

	return s.startStep(
		step,
		removeAfterDone,
		noNewLineAtStart,
	)
}

func (s Stepper) StartTemporaryStepWithoutNewLine(
	step string,
) stepper.Step {

	removeAfterDone := true
	noNewLineAtStart := true

	return s.startStep(
		step,
		removeAfterDone,
		noNewLineAtStart,
	)
}

func (s Stepper) StopCurrentStep() {

	if currentStep != nil {
		currentStep.Done()
		currentStep = nil
	}
}
