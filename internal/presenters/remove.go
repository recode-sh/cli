package presenters

import (
	"github.com/recode-sh/cli/internal/features"
)

type RemoveViewDataContent struct {
	Message string
}

type RemoveViewData struct {
	Error   *ViewableError
	Content RemoveViewDataContent
}

type RemoveViewer interface {
	View(RemoveViewData)
}

type RemovePresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               RemoveViewer
}

func NewRemovePresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer RemoveViewer,
) RemovePresenter {

	return RemovePresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (r RemovePresenter) PresentToView(response features.RemoveResponse) {
	viewData := RemoveViewData{}

	if response.Error == nil {
		devEnvName := response.Content.DevEnvName
		viewDataMessage := "The development environment \"" + devEnvName + "\" was removed."

		viewData.Content = RemoveViewDataContent{
			Message: viewDataMessage,
		}

		r.viewer.View(viewData)

		return
	}

	viewData.Error = r.viewableErrorBuilder.Build(response.Error)

	r.viewer.View(viewData)
}
