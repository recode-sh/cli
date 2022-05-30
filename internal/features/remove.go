package features

import (
	"github.com/recode-sh/recode/features"
)

type RemoveResponse struct {
	Error   error
	Content RemoveResponseContent
}

type RemoveResponseContent struct {
	DevEnvName string
}

type RemovePresenter interface {
	PresentToView(RemoveResponse)
}

type RemoveOutputHandler struct {
	presenter RemovePresenter
}

func NewRemoveOutputHandler(
	presenter RemovePresenter,
) RemoveOutputHandler {

	return RemoveOutputHandler{
		presenter: presenter,
	}
}

func (r RemoveOutputHandler) HandleOutput(output features.RemoveOutput) error {
	output.Stepper.StopCurrentStep()

	handleError := func(err error) error {
		r.presenter.PresentToView(RemoveResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	devEnv := output.Content.DevEnv

	r.presenter.PresentToView(RemoveResponse{
		Content: RemoveResponseContent{
			DevEnvName: devEnv.Name,
		},
	})

	return nil
}
