package presenters

import (
	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/features"
)

type StartViewData struct {
	Error   *ViewableError
	Content StartViewDataContent
}

type StartViewDataContent struct {
	ShowAsWarning bool
	Message       string
	Subtext       string
}

type StartViewer interface {
	View(StartViewData)
}

type StartPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               StartViewer
}

func NewStartPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer StartViewer,
) StartPresenter {

	return StartPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (s StartPresenter) PresentToView(response features.StartResponse) {
	viewData := StartViewData{}

	if response.Error == nil {
		devEnvName := response.Content.DevEnvName

		viewDataMessage := "The development environment \"" + devEnvName + "\" was started."
		viewDataSubtext := "Run `" + constants.Blue("ssh "+devEnvName) + "` (or use your code editor's integrated terminal)."

		devEnvAlreadyStarted := response.Content.DevEnvAlreadyStarted

		if devEnvAlreadyStarted {
			viewDataMessage = "The development environment \"" + devEnvName + "\" is already started."
			viewDataSubtext = ""
		}

		devEnvRebuilt := response.Content.DevEnvRebuilt

		if devEnvRebuilt {
			viewDataMessage = "The development environment \"" + devEnvName + "\" was rebuilt."
			viewDataSubtext = ""
		}

		viewData.Content = StartViewDataContent{
			ShowAsWarning: devEnvAlreadyStarted,
			Message:       viewDataMessage,
			Subtext:       viewDataSubtext,
		}

		s.viewer.View(viewData)

		return
	}

	viewData.Error = s.viewableErrorBuilder.Build(response.Error)

	s.viewer.View(viewData)
}
