package views

import (
	"io"
	"os"

	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/presenters"
)

//go:generate mockgen -destination=../mocks/views_displayer.go -package=mocks github.com/recode-sh/cli/internal/views Displayer
type Displayer interface {
	Display(w io.Writer, format string, args ...interface{})
}

type BaseView struct {
	Displayer Displayer
}

func NewBaseView(displayer Displayer) BaseView {
	return BaseView{
		Displayer: displayer,
	}
}

func (b BaseView) showErrorView(
	err *presenters.ViewableError,
	startWithNewLine bool,
) {

	bold := constants.Bold
	red := constants.Red

	if startWithNewLine {
		b.Displayer.Display(
			os.Stdout,
			"\n",
		)
	}

	b.Displayer.Display(
		os.Stdout,
		"%s %s\n\n%s\n\n",
		bold(red("Error!")),
		bold(err.Title),
		err.Message,
	)
}

func (b BaseView) ShowErrorView(err *presenters.ViewableError) {
	b.showErrorView(err, false)
}

func (b BaseView) ShowErrorViewWithStartingNewLine(err *presenters.ViewableError) {
	b.showErrorView(err, true)
}

func (b BaseView) ShowWarningView(warningText, subtext string) {
	bold := constants.Bold
	yellow := constants.Yellow

	if len(subtext) > 0 {
		b.Displayer.Display(
			os.Stdout,
			"%s %s\n\n%s\n\n",
			bold(yellow("Warning!")),
			bold(warningText),
			subtext,
		)

		return
	}

	b.Displayer.Display(
		os.Stdout,
		"%s %s\n\n",
		bold(yellow("Warning!")),
		bold(warningText),
	)
}

func (b BaseView) ShowSuccessView(successText, subtext string) {
	bold := constants.Bold
	green := constants.Green

	if len(subtext) > 0 {
		b.Displayer.Display(
			os.Stdout,
			"%s %s\n\n%s\n\n",
			bold(green("Success!")),
			bold(successText),
			subtext,
		)

		return
	}

	b.Displayer.Display(
		os.Stdout,
		"%s %s\n\n",
		bold(green("Success!")),
		bold(successText),
	)
}
