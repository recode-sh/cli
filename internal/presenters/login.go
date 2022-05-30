package presenters

import (
	"github.com/recode-sh/cli/internal/features"
)

type LoginViewDataContent struct {
	Message string
}

type LoginViewData struct {
	Error   *ViewableError
	Content LoginViewDataContent
}

type LoginViewer interface {
	View(LoginViewData)
}

type LoginPresenter struct {
	viewableErrorBuilder ViewableErrorBuilder
	viewer               LoginViewer
}

func NewLoginPresenter(
	viewableErrorBuilder ViewableErrorBuilder,
	viewer LoginViewer,
) LoginPresenter {

	return LoginPresenter{
		viewableErrorBuilder: viewableErrorBuilder,
		viewer:               viewer,
	}
}

func (l LoginPresenter) PresentToView(response features.LoginResponse) {
	viewData := LoginViewData{}

	if response.Error == nil {
		viewDataMessage := "Your GitHub account is now connected."

		viewData.Content = LoginViewDataContent{
			Message: viewDataMessage,
		}

		l.viewer.View(viewData)

		return
	}

	viewData.Error = l.viewableErrorBuilder.Build(response.Error)

	l.viewer.View(viewData)
}
