package vscode

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/recode-sh/cli/internal/system"
)

type ErrCLINotFound struct {
	VisitedPaths []string
}

func (ErrCLINotFound) Error() string {
	return "ErrCLINotFound"
}

type CLI struct{}

func (c CLI) Exec(arg ...string) (string, error) {
	CLIPath, err := c.LookupPath(runtime.GOOS)

	if err != nil {
		return "", err
	}

	cmd := exec.Command(CLIPath, arg...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		newLineRegExp := regexp.MustCompile(`\n+`)

		return "", fmt.Errorf(
			"Error while calling the Visual Studio Code CLI:\n\n%s\n\n%s",
			strings.TrimSpace(
				newLineRegExp.ReplaceAllLiteralString(stderr.String(), " "),
			),
			err.Error(),
		)
	}

	return stdout.String(), nil
}

func (c CLI) LookupPath(operatingSystem string) (string, error) {
	// First, we look for the 'code-insiders' command
	insidersCLIPath, err := exec.LookPath("code-insiders")

	if err == nil { // 'code-insiders' command exists
		return insidersCLIPath, nil
	}

	// If the 'code-insiders' command was not found, we look for the 'code' one
	CLIPath, err := exec.LookPath("code")

	if err == nil { // 'code' command exists
		return CLIPath, nil
	}

	// Finally, we fallback to default paths
	possibleCLIPaths := []string{}

	if operatingSystem == "darwin" { // macOS
		possibleCLIPaths = c.macOSPossibleCLIPaths()
	}

	if operatingSystem == "windows" {
		possibleCLIPaths = c.windowsPossibleCLIPaths()
	}

	if operatingSystem == "linux" {
		possibleCLIPaths = c.linuxPossibleCLIPaths()
	}

	for _, possibleCLIPath := range possibleCLIPaths {
		if system.PathExists(possibleCLIPath) {
			return possibleCLIPath, nil
		}
	}

	return "", ErrCLINotFound{
		VisitedPaths: possibleCLIPaths,
	}
}

func (c CLI) macOSPossibleCLIPaths() []string {
	rootApplicationsDir := fmt.Sprintf("%cApplications", os.PathSeparator) // /Applications

	// Order matter here.
	// We want the insiders version to be matched first.
	possiblePaths := []string{}

	possiblePaths = append(possiblePaths, filepath.Join(
		rootApplicationsDir,
		"Visual Studio Code - Insiders.app",
		"Contents",
		"Resources",
		"app",
		"bin",
		"code-insiders",
	))

	possiblePaths = append(possiblePaths, filepath.Join(
		rootApplicationsDir,
		"Visual Studio Code.app",
		"Contents",
		"Resources",
		"app",
		"bin",
		"code",
	))

	return possiblePaths
}

func (c CLI) windowsPossibleCLIPaths() []string {
	programFilesPath := os.Getenv("ProgramFiles")

	// Order matter here.
	// We want the insiders version to be matched first.

	// -- Insiders VSCode versions

	possiblePaths := []string{}

	possiblePaths = append(possiblePaths, filepath.Join(
		system.UserHomeDir(),
		"AppData",
		"Local",
		"Programs",
		"Microsoft VS Code Insiders",
		"bin",
		"code-insiders.cmd",
	))

	possiblePaths = append(possiblePaths, filepath.Join(
		programFilesPath,
		"Microsoft VS Code Insiders",
		"bin",
		"code-insiders.cmd",
	))

	// -- Regular VSCode versions

	possiblePaths = append(possiblePaths, filepath.Join(
		system.UserHomeDir(),
		"AppData",
		"Local",
		"Programs",
		"Microsoft VS Code",
		"bin",
		"code.cmd",
	))

	possiblePaths = append(possiblePaths, filepath.Join(
		programFilesPath,
		"Microsoft VS Code",
		"bin",
		"code.cmd",
	))

	return possiblePaths
}

func (c CLI) linuxPossibleCLIPaths() []string {
	// Order matter here.
	// We want the insiders version to be matched first.
	possiblePaths := []string{
		"/usr/bin/code-insiders",
		"/snap/bin/code-insiders",
		"/usr/share/code/bin/code-insiders",

		"/usr/bin/code",
		"/snap/bin/code",
		"/usr/share/code/bin/code",
	}

	return possiblePaths
}
