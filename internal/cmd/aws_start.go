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
	"os"

	"github.com/recode-sh/cli/internal/dependencies"
	"github.com/recode-sh/cli/internal/system"
	"github.com/recode-sh/recode/features"
	"github.com/spf13/cobra"
)

var awsStartInstanceType string
var awsStartDevEnvRebuildAsked bool
var awsStartForceDevEnvRebuild bool

// awsStartCmd represents the aws start command
var awsStartCmd = &cobra.Command{
	Use: "start <repository>",

	Short: "Start a development environment",

	Long: `Start a development environment for a specific GitHub repository.

If the passed repository doesn't contain an account name, your personal account is assumed.`,

	Example: `  recode aws start api
  recode aws start recode-sh/cli --instance-type m4.large`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		recodeViewableErrorBuilder := dependencies.ProvideRecodeViewableErrorBuilder()
		baseView := dependencies.ProvideBaseView()

		devEnvUserConfigResolver := dependencies.ProvideDevEnvUserConfigResolver()
		resolvedDevEnvUserConfig, err := devEnvUserConfigResolver.Resolve()

		if err != nil {
			baseView.ShowErrorViewWithStartingNewLine(
				recodeViewableErrorBuilder.Build(
					err,
				),
			)

			os.Exit(1)
		}

		repository := args[0]
		checkForRepositoryExistence := true

		repositoryResolver := dependencies.ProvideDevEnvRepositoryResolver()
		resolvedRepository, err := repositoryResolver.Resolve(
			repository,
			checkForRepositoryExistence,
		)

		if err != nil {
			baseView.ShowErrorViewWithStartingNewLine(
				recodeViewableErrorBuilder.Build(
					err,
				),
			)

			os.Exit(1)
		}

		awsStartInput := features.StartInput{
			InstanceType:             awsStartInstanceType,
			DevEnvRebuildAsked:       awsStartDevEnvRebuildAsked,
			ResolvedDevEnvUserConfig: *resolvedDevEnvUserConfig,
			ResolvedRepository:       *resolvedRepository,
			ForceDevEnvRevuild:       awsStartForceDevEnvRebuild,
			ConfirmDevEnvRebuild: func() (bool, error) {
				logger := system.NewLogger()
				return system.AskForConfirmation(
					logger,
					os.Stdin,
					"All your un-pushed work will be lost",
				)
			},
		}

		awsStart := dependencies.ProvideAWSStartFeature(
			awsRegion,
			awsProfile,
			awsCredentialsFilePath,
			awsConfigFilePath,
		)

		err = awsStart.Execute(awsStartInput)

		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	awsStartCmd.Flags().StringVar(
		&awsStartInstanceType,
		"instance-type",
		"t2.medium",
		"the instance type used by this development environment",
	)

	awsStartCmd.Flags().BoolVar(
		&awsStartDevEnvRebuildAsked,
		"rebuild",
		false,
		"rebuild the development environment",
	)

	awsStartCmd.Flags().BoolVar(
		&awsStartForceDevEnvRebuild,
		"force",
		false,
		"avoid rebuild confirmation",
	)

	awsCmd.AddCommand(awsStartCmd)
}
