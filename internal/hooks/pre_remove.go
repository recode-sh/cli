package hooks

import (
	"encoding/json"

	"github.com/recode-sh/cli/internal/config"
	cliEntities "github.com/recode-sh/cli/internal/entities"
	"github.com/recode-sh/cli/internal/interfaces"
	"github.com/recode-sh/recode/entities"
)

type PreRemove struct {
	sshConfig     interfaces.SSHConfigManager
	sshKeys       interfaces.SSHKeysManager
	sshKnownHosts interfaces.SSHKnownHostsManager
	userConfig    interfaces.UserConfigManager
	github        interfaces.GitHubManager
}

func NewPreRemove(
	sshConfig interfaces.SSHConfigManager,
	sshKeys interfaces.SSHKeysManager,
	sshKnownHosts interfaces.SSHKnownHostsManager,
	userConfig interfaces.UserConfigManager,
	github interfaces.GitHubManager,
) PreRemove {

	return PreRemove{
		sshConfig:     sshConfig,
		sshKeys:       sshKeys,
		sshKnownHosts: sshKnownHosts,
		userConfig:    userConfig,
		github:        github,
	}
}

func (p PreRemove) Run(
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
) error {

	err := p.sshKeys.RemovePEMIfExists(devEnv.GetSSHKeyPairName())

	if err != nil {
		return err
	}

	sshConfigHostKey := devEnv.Name
	err = p.sshConfig.RemoveHostIfExists(sshConfigHostKey)

	if err != nil {
		return err
	}

	sshHostname := devEnv.InstancePublicIPAddress
	err = p.sshKnownHosts.RemoveIfExists(sshHostname)

	if err != nil {
		return err
	}

	// User could remove dev env in creating state
	// (in case of error for example)
	if len(devEnv.AdditionalPropertiesJSON) == 0 {
		return nil
	}

	var devEnvAdditionalProperties *cliEntities.DevEnvAdditionalProperties
	err = json.Unmarshal(
		[]byte(devEnv.AdditionalPropertiesJSON),
		&devEnvAdditionalProperties,
	)

	if err != nil {
		return err
	}

	githubAccessToken := p.userConfig.GetString(
		config.UserConfigKeyGitHubAccessToken,
	)

	if devEnvAdditionalProperties.GitHubCreatedSSHKeyId != nil {
		err = p.github.RemoveSSHKey(
			githubAccessToken,
			*devEnvAdditionalProperties.GitHubCreatedSSHKeyId,
		)

		if err != nil && !p.github.IsNotFoundError(err) {
			return err
		}
	}

	if devEnvAdditionalProperties.GitHubCreatedGPGKeyId != nil {
		err = p.github.RemoveGPGKey(
			githubAccessToken,
			*devEnvAdditionalProperties.GitHubCreatedGPGKeyId,
		)

		if err != nil && !p.github.IsNotFoundError(err) {
			return err
		}
	}

	return nil
}
