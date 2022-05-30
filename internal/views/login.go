package views

import "github.com/recode-sh/cli/internal/presenters"

type LoginView struct {
	BaseView
}

func NewLoginView(baseView BaseView) LoginView {
	return LoginView{
		BaseView: baseView,
	}
}

func (l LoginView) View(data presenters.LoginViewData) {
	if data.Error == nil {
		l.ShowSuccessView(data.Content.Message, "")
		return
	}

	l.ShowErrorView(data.Error)
}
