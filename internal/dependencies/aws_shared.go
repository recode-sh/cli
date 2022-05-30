// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/google/wire"
	awsProviderConfig "github.com/recode-sh/aws-cloud-provider/config"
	awsProviderService "github.com/recode-sh/aws-cloud-provider/service"
	awsProviderUserConfig "github.com/recode-sh/aws-cloud-provider/userconfig"
	awsCLI "github.com/recode-sh/cli/internal/aws"
	"github.com/recode-sh/cli/internal/presenters"
	"github.com/recode-sh/cli/internal/system"
	"github.com/recode-sh/recode/entities"
)

var awsViewableErrorBuilder = wire.NewSet(
	wire.Bind(new(presenters.ViewableErrorBuilder), new(awsCLI.AWSViewableErrorBuilder)),
	awsCLI.NewAWSViewableErrorBuilder,
)

var awsServiceBuilderSet = wire.NewSet(
	wire.Bind(new(awsProviderUserConfig.ProfileLoader), new(awsProviderConfig.ProfileLoader)),
	awsProviderConfig.NewProfileLoader,

	wire.Bind(new(awsCLI.UserConfigFilesResolver), new(awsProviderUserConfig.FilesResolver)),
	awsProviderUserConfig.NewFilesResolver,

	wire.Bind(new(awsProviderUserConfig.EnvVarsGetter), new(system.EnvVars)),
	system.NewEnvVars,

	wire.Bind(new(awsCLI.UserConfigEnvVarsResolver), new(awsProviderUserConfig.EnvVarsResolver)),
	awsProviderUserConfig.NewEnvVarsResolver,

	wire.Bind(new(awsProviderService.UserConfigResolver), new(awsCLI.UserConfigLocalResolver)),
	awsCLI.NewUserConfigLocalResolver,

	wire.Bind(new(awsProviderService.UserConfigValidator), new(awsProviderConfig.UserConfigValidator)),
	awsProviderConfig.NewUserConfigValidator,

	wire.Bind(new(awsProviderService.UserConfigLoader), new(awsProviderConfig.UserConfigLoader)),
	awsProviderConfig.NewUserConfigLoader,

	wire.Bind(new(entities.CloudServiceBuilder), new(awsProviderService.Builder)),
	awsProviderService.NewBuilder,
)
