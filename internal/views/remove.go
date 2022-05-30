package views

import "github.com/recode-sh/cli/internal/presenters"

type RemoveView struct {
	BaseView
}

func NewRemoveView(baseView BaseView) RemoveView {
	return RemoveView{
		BaseView: baseView,
	}
}

func (r RemoveView) View(data presenters.RemoveViewData) {
	if data.Error == nil {
		r.ShowSuccessView(data.Content.Message, "")
		return
	}

	r.ShowErrorView(data.Error)
}
