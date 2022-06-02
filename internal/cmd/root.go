/*
Copyright © 2022 Jeremy Levy jje.levy@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/recode-sh/cli/internal/config"
	"github.com/recode-sh/cli/internal/dependencies"
	"github.com/recode-sh/cli/internal/exceptions"
	"github.com/recode-sh/cli/internal/system"
	"github.com/recode-sh/cli/internal/vscode"
	"github.com/recode-sh/recode/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "recode",

	Short: "Remote development environments defined as code",

	Long: `Recode - Remote development environments defined as code

To begin, run the command "recode login" to connect your GitHub account.	

From there, the most common workflow is:

  - recode <cloud_provider> start <repository>  : to start a development environment for a specific GitHub repository
  - recode <cloud_provider> stop <repository>   : to stop a development environment (without removing your data)
  - recode <cloud_provider> remove <repository> : to remove a development environment AND your data
	
<repository> may be relative to your personal GitHub account (eg: cli) or fully qualified (eg: my-organization/api).`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ensureUserIsLoggedIn(cmd)
	},

	TraverseChildren: true,

	Version: "v0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		ensureRecodeCLIRequirements,
		initializeRecodeCLIConfig,
		ensureGitHubAccessTokenValidity,
	)
}

func ensureRecodeCLIRequirements() {
	missingRequirements := []string{}

	vscodeCLI := vscode.CLI{}
	_, err := vscodeCLI.LookupPath(runtime.GOOS)

	if vscodeCLINotFoundErr, ok := err.(vscode.ErrCLINotFound); ok {
		missingRequirements = append(
			missingRequirements,
			fmt.Sprintf(
				"Visual Studio Code (looked in \"%s)",
				strings.Join(vscodeCLINotFoundErr.VisitedPaths, "\", \"")+"\"",
			),
		)
	}

	sshCommand := "ssh"
	_, err = exec.LookPath(sshCommand)

	if err != nil {
		missingRequirements = append(
			missingRequirements,
			fmt.Sprintf(
				"OpenSSH client (looked for an \"%s\" command available)",
				sshCommand,
			),
		)
	}

	if len(missingRequirements) > 0 {
		recodeViewableErrorBuilder := dependencies.ProvideRecodeViewableErrorBuilder()
		baseView := dependencies.ProvideBaseView()

		missingRequirementsErr := exceptions.ErrMissingRequirements{
			MissingRequirements: missingRequirements,
		}

		baseView.ShowErrorViewWithStartingNewLine(
			recodeViewableErrorBuilder.Build(
				missingRequirementsErr,
			),
		)

		os.Exit(1)
	}
}

func initializeRecodeCLIConfig() {
	configDir := system.UserConfigDir()
	configDirPerms := fs.FileMode(0700)

	// Ensure configuration dir exists
	err := os.MkdirAll(
		configDir,
		configDirPerms,
	)
	cobra.CheckErr(err)

	configFilePath := system.UserConfigFilePath()
	configFilePerms := fs.FileMode(0600)

	// Ensure configuration file exists
	f, err := os.OpenFile(
		configFilePath,
		os.O_CREATE,
		configFilePerms,
	)
	cobra.CheckErr(err)
	defer f.Close()

	viper.SetConfigFile(configFilePath)
	cobra.CheckErr(viper.ReadInConfig())
}

// ensureGitHubAccessTokenValidity ensures that
// the github access token has not been
// revoked by user
func ensureGitHubAccessTokenValidity() {
	userConfig := config.NewUserConfig()
	userIsLoggedIn := userConfig.GetBool(config.UserConfigKeyUserIsLoggedIn)

	if !userIsLoggedIn {
		return
	}

	gitHubService := github.NewService()

	githubUser, err := gitHubService.GetAuthenticatedUser(
		userConfig.GetString(
			config.UserConfigKeyGitHubAccessToken,
		),
	)

	if err != nil &&
		gitHubService.IsInvalidAccessTokenError(err) { // User has revoked access token

		userIsLoggedIn = false

		userConfig.Set(
			config.UserConfigKeyUserIsLoggedIn,
			userIsLoggedIn,
		)

		// Error is swallowed here to
		// not confuse user with unexpected error
		_ = userConfig.WriteConfig()
	}

	if err == nil {
		// Update config with updated values from GitHub
		userConfig.PopulateFromGitHubUser(githubUser)

		// Error is swallowed here to
		// not confuse user with unexpected error
		_ = userConfig.WriteConfig()
	}
}

func ensureUserIsLoggedIn(cmd *cobra.Command) {
	userConfig := config.NewUserConfig()
	userIsLoggedIn := userConfig.GetBool(config.UserConfigKeyUserIsLoggedIn)

	if !userIsLoggedIn && cmd != loginCmd {
		recodeViewableErrorBuilder := dependencies.ProvideRecodeViewableErrorBuilder()
		baseView := dependencies.ProvideBaseView()

		baseView.ShowErrorViewWithStartingNewLine(
			recodeViewableErrorBuilder.Build(
				exceptions.ErrUserNotLoggedIn,
			),
		)

		os.Exit(1)
	}
}
