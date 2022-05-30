package views

import "github.com/recode-sh/cli/internal/presenters"

type UninstallView struct {
	BaseView
}

func NewUninstallView(baseView BaseView) UninstallView {
	return UninstallView{
		BaseView: baseView,
	}
}

func (u UninstallView) View(data presenters.UninstallViewData) {
	if data.Error == nil {
		if data.Content.ShowAsWarning {
			u.ShowWarningView(
				data.Content.Message,
				data.Content.Subtext,
			)
			return
		}

		u.ShowSuccessView(
			data.Content.Message,
			data.Content.Subtext,
		)
		return
	}

	u.ShowErrorView(data.Error)
}
