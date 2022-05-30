// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/google/wire"
	awsProviderUserConfig "github.com/recode-sh/aws-cloud-provider/userconfig"
	"github.com/recode-sh/cli/internal/agent"
	awsCLI "github.com/recode-sh/cli/internal/aws"
	featuresCLI "github.com/recode-sh/cli/internal/features"
	"github.com/recode-sh/cli/internal/presenters"
	"github.com/recode-sh/cli/internal/views"
	"github.com/recode-sh/recode/features"
)

func ProvideAWSStartFeature(region, profile, credentialsFilePath, configFilePath string) features.StartFeature {
	return provideAWSStartFeature(
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

func provideAWSStartFeature(
	userConfigEnvVarsResolverOpts awsProviderUserConfig.EnvVarsResolverOpts,
	userConfigFilesResolverOpts awsProviderUserConfig.FilesResolverOpts,
	userConfigLocalResolverOpts awsCLI.UserConfigLocalResolverOpts,
) features.StartFeature {
	panic(
		wire.Build(
			viewSet,
			awsServiceBuilderSet,
			awsViewableErrorBuilder,

			userConfigManagerSet,

			wire.Bind(new(agent.ClientBuilder), new(agent.DefaultClientBuilder)),
			agent.NewDefaultClientBuilder,

			githubManagerSet,

			loggerSet,

			sshConfigManagerSet,

			sshKnownHostsManagerSet,

			sshKeysManagerSet,

			vscodeProcessManagerSet,

			vscodeExtensionsManagerSet,

			stepperSet,

			wire.Bind(new(features.StartOutputHandler), new(featuresCLI.StartOutputHandler)),
			featuresCLI.NewStartOutputHandler,

			wire.Bind(new(featuresCLI.StartPresenter), new(presenters.StartPresenter)),
			presenters.NewStartPresenter,

			wire.Bind(new(presenters.StartViewer), new(views.StartView)),
			views.NewStartView,

			features.NewStartFeature,
		),
	)
}
