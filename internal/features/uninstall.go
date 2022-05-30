package features

import (
	"os"

	"github.com/recode-sh/cli/internal/system"
	"github.com/recode-sh/recode/features"
)

type UninstallResponse struct {
	Error   error
	Content UninstallResponseContent
}

type UninstallResponseContent struct {
	RecodeAlreadyUninstalled  bool
	SuccessMessage            string
	AlreadyUninstalledMessage string
	RecodeExecutablePath      string
	RecodeConfigDirPath       string
}

type UninstallPresenter interface {
	PresentToView(UninstallResponse)
}

type UninstallOutputHandler struct {
	presenter UninstallPresenter
}

func NewUninstallOutputHandler(
	presenter UninstallPresenter,
) UninstallOutputHandler {

	return UninstallOutputHandler{
		presenter: presenter,
	}
}

func (u UninstallOutputHandler) HandleOutput(output features.UninstallOutput) error {
	output.Stepper.StopCurrentStep()

	handleError := func(err error) error {
		u.presenter.PresentToView(UninstallResponse{
			Error: err,
		})

		return err
	}

	if output.Error != nil {
		return handleError(output.Error)
	}

	recodeExecutablePath, err := os.Executable()

	if err != nil {
		return handleError(err)
	}

	recodeConfigDirPath := system.UserConfigDir()

	u.presenter.PresentToView(UninstallResponse{
		Content: UninstallResponseContent{
			RecodeAlreadyUninstalled:  output.Content.RecodeAlreadyUninstalled,
			SuccessMessage:            output.Content.SuccessMessage,
			AlreadyUninstalledMessage: output.Content.AlreadyUninstalledMessage,
			RecodeExecutablePath:      recodeExecutablePath,
			RecodeConfigDirPath:       recodeConfigDirPath,
		},
	})

	return nil
}
