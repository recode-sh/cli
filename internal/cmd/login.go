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
	"github.com/recode-sh/cli/internal/features"
	"github.com/spf13/cobra"
)

// loginCmd represents the "recode login" command
var loginCmd = &cobra.Command{
	Use: "login",

	Short: "Connect a GitHub account to use with Recode",

	Long: `Connect a GitHub account to use with Recode.

Recode requires the following permissions:

  - "Public SSH keys" and "Repositories" to let you access your repositories from your development environments

  - "GPG Keys" and "Personal user data" to configure Git and sign your commits (verified badge)

All your data (including the OAuth access token) will only be stored locally.`,

	Example: "  recode login",

	Run: func(cmd *cobra.Command, args []string) {
		loginInput := features.LoginInput{}

		login := dependencies.ProvideLoginFeature()

		err := login.Execute(loginInput)

		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
