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

func ProvideAWSStopFeature(region, profile, credentialsFilePath, configFilePath string) features.StopFeature {
	return provideAWSStopFeature(
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

func provideAWSStopFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.StopFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			sshKnownHostsManagerSet,

			stepperSet,

			wire.Bind(new(features.StopOutputHandler), new(featuresCLI.StopOutputHandler)),
			featuresCLI.NewStopOutputHandler,

			wire.Bind(new(featuresCLI.StopPresenter), new(presenters.StopPresenter)),
			presenters.NewStopPresenter,

			wire.Bind(new(presenters.StopViewer), new(views.StopView)),
			views.NewStopView,

			features.NewStopFeature,
		),
	)
}
