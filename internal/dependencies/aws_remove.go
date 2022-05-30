// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/google/wire"
	awsProviderUserConfig "github.com/recode-sh/aws-cloud-provider/userconfig"
	awsCLI "github.com/recode-sh/cli/internal/aws"
	featuresCLI "github.com/recode-sh/cli/internal/features"
	"github.com/recode-sh/cli/internal/presenters"
	"github.com/recode-sh/cli/internal/views"
	"github.com/recode-sh/recode/features"
)

func ProvideAWSRemoveFeature(region, profile, credentialsFilePath, configFilePath string) features.RemoveFeature {
	return provideAWSRemoveFeature(
		awsProviderUserConfig.EnvVarsResolverOpts{
			Region: region,
		},

		awsProviderUserConfig.FilesResolverOpts{
			Region:              region,
			Profile:             profile,
			CredentialsFilePath: credentialsFilePath,
			ConfigFilePath:      configFilePath,
		},

		awsCLI.UserConfigLocalResolverOpts{
			Profile: profile,
		},
	)
}

func provideAWSRemoveFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.RemoveFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			stepperSet,

			wire.Bind(new(features.RemoveOutputHandler), new(featuresCLI.RemoveOutputHandler)),
			featuresCLI.NewRemoveOutputHandler,

			wire.Bind(new(featuresCLI.RemovePresenter), new(presenters.RemovePresenter)),
			presenters.NewRemovePresenter,

			wire.Bind(new(presenters.RemoveViewer), new(views.RemoveView)),
			views.NewRemoveView,

			features.NewRemoveFeature,
		),
	)
}
