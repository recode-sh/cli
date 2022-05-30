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

// awsUninstallCmd represents the aws uninstall command
var awsUninstallCmd = &cobra.Command{
	Use: "uninstall",

	Short: "Uninstall Recode from your AWS account",

	Long: `Uninstall Recode from your AWS account.

All your development environments must be removed before running this command.`,

	Example: "  recode aws uninstall",

	Run: func(cmd *cobra.Command, args []string) {

		awsUninstallInput := features.UninstallInput{
			SuccessMessage:            "Recode has been uninstalled from this region on this AWS account.",
			AlreadyUninstalledMessage: "Recode is already uninstalled in this region on this AWS account.",
		}

		awsUninstall := dependencies.ProvideAWSUninstallFeature(
			awsRegion,
			awsProfile,
			awsCredentialsFilePath,
			awsConfigFilePath,
		)

		err := awsUninstall.Execute(awsUninstallInput)

		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	awsCmd.AddCommand(awsUninstallCmd)
}
