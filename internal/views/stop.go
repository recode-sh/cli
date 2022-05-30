package views

import "github.com/recode-sh/cli/internal/presenters"

type StopView struct {
	BaseView
}

func NewStopView(baseView BaseView) StopView {
	return StopView{
		BaseView: baseView,
	}
}

func (s StopView) View(data presenters.StopViewData) {
	if data.Error == nil {
		if data.Content.ShowAsWarning {
			s.ShowWarningView(data.Content.Message, "")
			return
		}

		s.ShowSuccessView(data.Content.Message, "")
		return
	}

	s.ShowErrorView(data.Error)
}
