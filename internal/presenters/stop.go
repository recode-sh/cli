package presenters

import (
	"github.com/recode-sh/cli/internal/features"
)

type StopViewDataContent struct {
	ShowAsWarning bool
	Message       string
}

type StopViewData struct {
	Error   *ViewableError
	Content StopViewDataContent
}

type StopViewer interface {
	View(StopViewData)
}

type StopPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               StopViewer
}

func NewStopPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer StopViewer,
) StopPresenter {

	return StopPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (s StopPresenter) PresentToView(response features.StopResponse) {
	viewData := StopViewData{}

	if response.Error == nil {
		devEnvName := response.Content.DevEnvName
		devEnvAlreadyStopped := response.Content.DevEnvAlreadyStopped

		viewDataMessage := "The development environment \"" + devEnvName + "\" was stopped."

		if devEnvAlreadyStopped {
			viewDataMessage = "The development environment \"" + devEnvName + "\" is already stopped. Nothing to do."
		}

		viewData.Content = StopViewDataContent{
			ShowAsWarning: devEnvAlreadyStopped,
			Message:       viewDataMessage,
		}

		s.viewer.View(viewData)

		return
	}

	viewData.Error = s.viewableErrorBuilder.Build(response.Error)

	s.viewer.View(viewData)
}
