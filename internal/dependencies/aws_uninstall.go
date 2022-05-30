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

func ProvideAWSUninstallFeature(region, profile, credentialsFilePath, configFilePath string) features.UninstallFeature {
	return provideAWSUninstallFeature(
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

func provideAWSUninstallFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.UninstallFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			stepperSet,

			wire.Bind(new(features.UninstallOutputHandler), new(featuresCLI.UninstallOutputHandler)),
			featuresCLI.NewUninstallOutputHandler,

			wire.Bind(new(featuresCLI.UninstallPresenter), new(presenters.UninstallPresenter)),
			presenters.NewUninstallPresenter,

			wire.Bind(new(presenters.UninstallViewer), new(views.UninstallView)),
			views.NewUninstallView,

			features.NewUninstallFeature,
		),
	)
}
