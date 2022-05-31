package entities

import (
	"fmt"

	"github.com/recode-sh/cli/internal/config"
	"github.com/recode-sh/cli/internal/constants"
	"github.com/recode-sh/cli/internal/interfaces"
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/github"
)

type DevEnvUserConfigResolver struct {
	logger     interfaces.Logger
	userConfig interfaces.UserConfigManager
	github     interfaces.GitHubManager
}

func NewDevEnvUserConfigResolver(
	logger interfaces.Logger,
	userConfig interfaces.UserConfigManager,
	github interfaces.GitHubManager,
) DevEnvUserConfigResolver {

	return DevEnvUserConfigResolver{
		logger:     logger,
		userConfig: userConfig,
		github:     github,
	}
}

func (d DevEnvUserConfigResolver) Resolve() (
	*entities.ResolvedDevEnvUserConfig,
	error,
) {

	githubAccessToken := d.userConfig.GetString(
		config.UserConfigKeyGitHubAccessToken,
	)

	devEnvUserConfigRepoOwner := d.userConfig.GetString(
		config.UserConfigKeyGitHubUsername,
	)

	userHasDevEnvUserConfigRepo, err := d.github.DoesRepositoryExist(
		githubAccessToken,
		devEnvUserConfigRepoOwner,
		entities.DevEnvUserConfigRepoName,
	)

	if err != nil {
		return nil, err
	}

	if !userHasDevEnvUserConfigRepo {
		yellow := constants.Blue

		d.logger.Log(
			"\n%s No repository \"%s\" found in GitHub account \"%s\". The repository \"%s\" will be used instead.",
			yellow("[Info]"),
			entities.DevEnvUserConfigRepoName,
			devEnvUserConfigRepoOwner,
			entities.DevEnvUserConfigDefaultRepoOwner+"/"+entities.DevEnvUserConfigRepoName,
		)

		devEnvUserConfigRepoOwner = entities.DevEnvUserConfigDefaultRepoOwner
	}

	_, err = d.github.GetFileContentFromRepository(
		githubAccessToken,
		devEnvUserConfigRepoOwner,
		entities.DevEnvUserConfigRepoName,
		entities.DevEnvUserConfigDockerfileFileName,
	)

	if err != nil && d.github.IsNotFoundError(err) {
		return nil, entities.ErrInvalidDevEnvUserConfig{
			RepoOwner: devEnvUserConfigRepoOwner,
			Reason: fmt.Sprintf(
				"Your repository must contain a file named \"%s\".",
				entities.DevEnvUserConfigDockerfileFileName,
			),
		}
	}

	if err != nil {
		return nil, err
	}

	return &entities.ResolvedDevEnvUserConfig{
		RepoOwner: devEnvUserConfigRepoOwner,
		RepoName:  entities.DevEnvUserConfigRepoName,

		RepoGitURL: github.BuildGitURL(
			devEnvUserConfigRepoOwner,
			entities.DevEnvUserConfigRepoName,
		),

		RepoGitHTTPURL: github.BuildGitHTTPURL(
			devEnvUserConfigRepoOwner,
			entities.DevEnvUserConfigRepoName,
		),
	}, nil
}
