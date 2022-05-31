package presenters

import (
	"fmt"

	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/features"
)

type UninstallViewDataContent struct {
	ShowAsWarning bool
	Message       string
	Subtext       string
}

type UninstallViewData struct {
	Error   *ViewableError
	Content UninstallViewDataContent
}

type UninstallViewer interface {
	View(UninstallViewData)
}

type UninstallPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               UninstallViewer
}

func NewUninstallPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer UninstallViewer,
) UninstallPresenter {

	return UninstallPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (u UninstallPresenter) PresentToView(response features.UninstallResponse) {
	viewData := UninstallViewData{}

	if response.Error == nil {
		bold := constants.Bold

		recodeAlreadyUninstalled := response.Content.RecodeAlreadyUninstalled

		viewDataMessage := response.Content.SuccessMessage
		viewDataSubtext := fmt.Sprintf(
			"If you want to remove Recode entirely:\n\n"+
				"  - Remove the Recode CLI (located at %s)\n\n"+
				"  - Remove the Recode configuration (located at %s)\n\n"+
				"  - Unauthorize the Recode application on GitHub by going to: %s",
			bold(response.Content.RecodeExecutablePath),
			bold(response.Content.RecodeConfigDirPath),
			bold("https://github.com/settings/applications"),
		)

		if recodeAlreadyUninstalled {
			viewDataMessage = response.Content.AlreadyUninstalledMessage
		}

		viewData.Content = UninstallViewDataContent{
			ShowAsWarning: recodeAlreadyUninstalled,
			Message:       viewDataMessage,
			Subtext:       viewDataSubtext,
		}

		u.viewer.View(viewData)

		return
	}

	viewData.Error = u.viewableErrorBuilder.Build(response.Error)
	u.viewer.View(viewData)
}
