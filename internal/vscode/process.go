package vscode

type Process struct{}

func NewProcess() Process {
	return Process{}
}

func (p Process) OpenOnRemote(hostKey, pathToOpen string) (string, error) {
	c := CLI{}

	return c.Exec(
		"--new-window",
		"--skip-release-notes",
		"--skip-welcome",
		"--skip-add-to-recently-opened",
		"--disable-workspace-trust",
		"--remote",
		"ssh-remote+"+hostKey,
		pathToOpen,
	)
}
