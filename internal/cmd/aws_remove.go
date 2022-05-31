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

var awsRemoveForceDevEnvRemove bool

// awsRemoveCmd represents the aws remove command
var awsRemoveCmd = &cobra.Command{
	Use: "remove <repository>",

	Short: "Remove a development environment",

	Long: `Remove an existing development environment.

The development environment will be PERMANENTLY removed along with ALL your data.
	
There is no going back, so please be sure to save your work before running this command.`,

	Example: `  recode aws remove api
  recode aws remove recode-sh/cli --force`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {

		recodeViewableErrorBuilder := dependencies.ProvideRecodeViewableErrorBuilder()
		baseView := dependencies.ProvideBaseView()

		repository := args[0]
		checkForRepositoryExistence := false

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

		awsRemoveInput := features.RemoveInput{
			ResolvedRepository: *resolvedRepository,
			PreRemoveHook:      dependencies.ProvidePreRemoveHook(),
			ForceRemove:        awsRemoveForceDevEnvRemove,
			ConfirmRemove: func() (bool, error) {
				logger := system.NewLogger()
				return system.AskForConfirmation(
					logger,
					os.Stdin,
					"All your un-pushed work will be lost.",
				)
			},
		}

		awsRemove := dependencies.ProvideAWSRemoveFeature(
			awsRegion,
			awsProfile,
			awsCredentialsFilePath,
			awsConfigFilePath,
		)

		err = awsRemove.Execute(awsRemoveInput)

		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	awsRemoveCmd.Flags().BoolVar(
		&awsRemoveForceDevEnvRemove,
		"force",
		false,
		"avoid remove confirmation",
	)

	awsCmd.AddCommand(awsRemoveCmd)
}
