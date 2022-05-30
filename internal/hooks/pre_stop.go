package hooks

import (
	"github.com/recode-sh/cli/internal/interfaces"
	"github.com/recode-sh/recode/entities"
)

type PreStop struct {
	sshKnownHosts interfaces.SSHKnownHostsManager
}

func NewPreStop(
	sshKnownHosts interfaces.SSHKnownHostsManager,
) PreStop {

	return PreStop{
		sshKnownHosts: sshKnownHosts,
	}
}

func (p PreStop) Run(
	cloudService entities.CloudService,
	recodeConfig *entities.Config,
	cluster *entities.Cluster,
	devEnv *entities.DevEnv,
) error {

	instanceIPAddress := devEnv.InstancePublicIPAddress

	return p.sshKnownHosts.RemoveIfExists(instanceIPAddress)
}
