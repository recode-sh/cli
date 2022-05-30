package entities

import (
	// "github.com/fatih/color"
	// "github.com/google/go-github/v43/github"
	"github.com/recode-sh/cli/internal/config"
	"github.com/recode-sh/cli/internal/interfaces"
	"github.com/recode-sh/recode/entities"
	"github.com/recode-sh/recode/github"
)

type DevEnvRepositoryResolver struct {
	logger     interfaces.Logger
	userConfig interfaces.UserConfigManager
	github     interfaces.GitHubManager
}

func NewDevEnvRepositoryResolver(
	logger interfaces.Logger,
	userConfig interfaces.UserConfigManager,
	github interfaces.GitHubManager,
) DevEnvRepositoryResolver {

	return DevEnvRepositoryResolver{
		logger:     logger,
		userConfig: userConfig,
		github:     github,
	}
}

func (d DevEnvRepositoryResolver) Resolve(
	repositoryName string,
	checkForRepositoryExistence bool,
) (*entities.ResolvedDevEnvRepository, error) {

	githubAccessToken := d.userConfig.GetString(
		config.UserConfigKeyGitHubAccessToken,
	)

	githubUsername := d.userConfig.GetString(
		config.UserConfigKeyGitHubUsername,
	)

	parsedRepoName, err := github.ParseRepositoryName(
		repositoryName,
		githubUsername,
	)

	if err != nil {
		// If repository name is invalid, we are sure
		// that the repository doesn't exist.
		return nil, entities.ErrDevEnvRepositoryNotFound{
			RepoOwner: githubUsername,
			RepoName:  repositoryName,
		}
	}

	if checkForRepositoryExistence {
		repoExists, err := d.github.DoesRepositoryExist(
			githubAccessToken,
			parsedRepoName.Owner,
			parsedRepoName.Name,
		)

		if err != nil {
			return nil, err
		}

		// if !repoExists && parsedRepoName.Owner != githubUsername {
		if !repoExists {
			return nil, entities.ErrDevEnvRepositoryNotFound{
				RepoOwner: parsedRepoName.Owner,
				RepoName:  parsedRepoName.Name,
			}
		}
	}

	// if !repoExists {
	// 	bold := color.New(color.Bold).SprintFunc()

	// 	d.logger.Log(
	// 		"\n%s "+bold("Repository \"%s\" not found. Creating now..."),
	// 		bold(color.YellowString("Warning!")),
	// 		parsedRepoName.Name,
	// 	)

	// 	// Means that we want the repository to be created
	// 	// in the logged user personal account. See GitHub SDK docs.
	// 	createdRepoOrganization := ""

	// 	createdRepoIsPrivate := true
	// 	createdRepoProps := &github.Repository{
	// 		Name:    &parsedRepoName.Name,
	// 		Private: &createdRepoIsPrivate,
	// 	}

	// 	_, err := d.github.CreateRepository(
	// 		githubAccessToken,
	// 		createdRepoOrganization,
	// 		createdRepoProps,
	// 	)

	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return &entities.ResolvedDevEnvRepository{
		Owner:         parsedRepoName.Owner,
		ExplicitOwner: parsedRepoName.ExplicitOwner,

		Name: parsedRepoName.Name,

		GitURL: github.BuildGitURL(
			parsedRepoName.Owner,
			parsedRepoName.Name,
		),

		GitHTTPURL: github.BuildGitHTTPURL(
			parsedRepoName.Owner,
			parsedRepoName.Name,
		),
	}, nil
}
