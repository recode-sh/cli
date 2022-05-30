// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/google/wire"
	"github.com/recode-sh/cli/internal/entities"
)

func ProvideDevEnvUserConfigResolver() entities.DevEnvUserConfigResolver {
	panic(
		wire.Build(
			loggerSet,

			userConfigManagerSet,

			githubManagerSet,

			entities.NewDevEnvUserConfigResolver,
		),
	)
}

func ProvideDevEnvRepositoryResolver() entities.DevEnvRepositoryResolver {
	panic(
		wire.Build(
			loggerSet,

			userConfigManagerSet,

			githubManagerSet,

			entities.NewDevEnvRepositoryResolver,
		),
	)
}
