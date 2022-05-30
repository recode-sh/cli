package interfaces

type VSCodeProcessManager interface {
	OpenOnRemote(hostKey, pathToOpen string) (cmdOutput string, cmdError error)
}

type VSCodeExtensionsManager interface {
	Install(extensionName string) (cmdOutput string, cmdError error)
}
