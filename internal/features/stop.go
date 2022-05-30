package features

import (
	"github.com/recode-sh/cli/internal/interfaces"
	"github.com/recode-sh/recode/features"
)

type StopResponse struct {
	Error   error
	Content StopResponseContent
}

type StopResponseContent struct {
	DevEnvName           string
	DevEnvAlreadyStopped bool
}

type StopPresenter interface {
	PresentToView(StopResponse)
}

type StopOutputHandler struct {
	presenter     StopPresenter
	sshKnownHosts interfaces.SSHKnownHostsManager
}

func NewStopOutputHandler(
	presenter StopPresenter,
	sshKnownHosts interfaces.SSHKnownHostsManager,
) StopOutputHandler {

	return StopOutputHandler{
		presenter:     presenter,
		sshKnownHosts: sshKnownHosts,
	}
}

func (s StopOutputHandler) HandleOutput(output features.StopOutput) error {
	output.Stepper.StopCurrentStep()

	handleError := func(err error) error {
		s.presenter.PresentToView(StopResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	devEnv := output.Content.DevEnv
	devEnvAlreadyStopped := output.Content.DevEnvAlreadyStopped

	if !devEnvAlreadyStopped {
		err := output.Content.SetDevEnvAsStopped()

		if err != nil {
			return handleError(err)
		}
	}

	s.presenter.PresentToView(StopResponse{
		Content: StopResponseContent{
			DevEnvName:           devEnv.Name,
			DevEnvAlreadyStopped: devEnvAlreadyStopped,
		},
	})

	return nil
}
