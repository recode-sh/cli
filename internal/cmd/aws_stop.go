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
	"github.com/recode-sh/recode/features"
	"github.com/spf13/cobra"
)

// awsStopCmd represents the aws stop command
var awsStopCmd = &cobra.Command{
	Use: "stop (<repository_name>|<account_name/repository_name>)",

	Short: "Stop a development environment",

	Long: `Stop an existing development environment.

The development environment will be stopped but your data will be conserved.

You may still incur charges for the storage used. If you don't plan to use this development environment again, use the remove command instead.`,

	Example: `  recode aws stop api
  recode aws stop recode-sh/cli`,

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

		awsStopInput := features.StopInput{
			ResolvedRepository: *resolvedRepository,
			PreStopHook:        dependencies.ProvidePreStopHook(),
		}

		awsStop := dependencies.ProvideAWSStopFeature(
			awsRegion,
			awsProfile,
			awsCredentialsFilePath,
			awsConfigFilePath,
		)

		err = awsStop.Execute(awsStopInput)

		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	awsCmd.AddCommand(awsStopCmd)
}
