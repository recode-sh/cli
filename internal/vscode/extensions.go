package vscode

type Extensions struct{}

func NewExtensions() Extensions {
	return Extensions{}
}

func (e Extensions) Install(extensionName string) (string, error) {
	c := CLI{}

	return c.Exec(
		"--install-extension",
		extensionName,
		"--force",
	)
}
