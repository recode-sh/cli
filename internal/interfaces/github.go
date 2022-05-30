package interfaces

import (
	"github.com/google/go-github/v43/github"
	cliGitHub "github.com/recode-sh/recode/github"
)

type GitHubManager interface {
	GetAuthenticatedUser(accessToken string) (*cliGitHub.AuthenticatedUser, error)

	CreateRepository(accessToken string, organization string, properties *github.Repository) (*github.Repository, error)
	DoesRepositoryExist(accessToken, repositoryOwner, repositoryName string) (bool, error)
	GetFileContentFromRepository(accessToken, repositoryOwner, repositoryName, filePath string) (string, error)

	CreateSSHKey(accessToken string, keyPairName string, publicKeyContent string) (*github.Key, error)
	RemoveSSHKey(accessToken string, sshKeyID int64) error

	CreateGPGKey(accessToken string, publicKeyContent string) (*github.GPGKey, error)
	RemoveGPGKey(accessToken string, gpgKeyID int64) error

	IsNotFoundError(err error) bool
}
