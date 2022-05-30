package views

import "github.com/recode-sh/cli/internal/presenters"

type StartView struct {
	BaseView
}

func NewStartView(baseView BaseView) StartView {
	return StartView{
		BaseView: baseView,
	}
}

func (s StartView) View(data presenters.StartViewData) {
	if data.Error == nil {
		if data.Content.ShowAsWarning {
			s.ShowWarningView(
				data.Content.Message,
				data.Content.Subtext,
			)
			return
		}

		s.ShowSuccessView(
			data.Content.Message,
			data.Content.Subtext,
		)
		return
	}

	s.ShowErrorView(data.Error)
}
