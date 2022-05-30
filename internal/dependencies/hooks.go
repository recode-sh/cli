// go:build wireinject
//go:build wireinject
// +build wireinject

package dependencies

import (
	"github.com/google/wire"
	"github.com/recode-sh/cli/internal/hooks"
)

func ProvidePreRemoveHook() hooks.PreRemove {
	panic(
		wire.Build(
			sshConfigManagerSet,

			sshKnownHostsManagerSet,

			sshKeysManagerSet,

			userConfigManagerSet,

			githubManagerSet,

			hooks.NewPreRemove,
		),
	)
}

func ProvidePreStopHook() hooks.PreStop {
	panic(
		wire.Build(

			sshKnownHostsManagerSet,

			hooks.NewPreStop,
		),
	)
}
