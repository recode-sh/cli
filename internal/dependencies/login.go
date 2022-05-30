// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/google/wire"
	"github.com/recode-sh/cli/internal/features"
	"github.com/recode-sh/cli/internal/presenters"
	"github.com/recode-sh/cli/internal/views"
)

func ProvideLoginFeature() features.LoginFeature {
	panic(
		wire.Build(
			viewSet,
			recodeViewableErrorBuilder,

			loggerSet,

			browserManagerSet,

			userConfigManagerSet,

			sleeperSet,

			githubManagerSet,

			wire.Bind(new(features.LoginPresenter), new(presenters.LoginPresenter)),
			presenters.NewLoginPresenter,

			wire.Bind(new(presenters.LoginViewer), new(views.LoginView)),
			views.NewLoginView,

			features.NewLoginFeature,
		),
	)
}
