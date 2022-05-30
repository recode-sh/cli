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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

var awsProfile string
var awsRegion string

var awsCredentialsFilePath string
var awsConfigFilePath string

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use: "aws",

	Short: "Use Recode on Amazon Web Services",

	Long: `Use Recode on Amazon Web Services.
	
To begin, create your first development environment using the command:

  recode aws start <repository>

Once started, you could stop it at any time, to save costs, using the command: 

  recode aws stop <repository>

If you don't plan to use this development environment again, you could remove it using the command:
	
  recode aws remove <repository>`,

	Example: `  recode aws start recode-sh/api --instance-type m4.large 
  recode aws stop recode-sh/api
  recode aws remove recode-sh/api`,
}

func init() {
	awsCmd.Flags().StringVar(
		&awsProfile,
		"profile",
		"",
		"the configuration profile to use to access your AWS account",
	)

	awsCmd.Flags().StringVar(
		&awsRegion,
		"region",
		"",
		"the region to use to access your AWS account",
	)

	awsCredentialsFilePath = config.DefaultSharedCredentialsFilename()
	awsConfigFilePath = config.DefaultSharedConfigFilename()

	rootCmd.AddCommand(awsCmd)
}
