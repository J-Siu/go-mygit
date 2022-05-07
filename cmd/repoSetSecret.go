/*
Copyright © 2022 John, Sing Dao, Siu <john.sd.siu@gmail.com>

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
	"os"

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-helper"
	"github.com/J-Siu/go-mygit/lib"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var repoSetSecretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"s", "secret"},
	Short:   "set action secrets",
	Run: func(cmd *cobra.Command, args []string) {
		for _, remote := range Conf.MergedRemotes {
			if remote.Vendor != "github" {
				fmt.Printf("%s(%s) action secret not supported.\n", remote.Name, remote.Vendor)
			} else {
				// "GET" public key
				helper.Report("", remote.Name, false, false)
				var pubkey gitapi.RepoPublicKey
				gitApi := lib.GitApiFromRemote(&remote, &pubkey, "")
				gitApi.EndpointReposSecretsPubkey()
				success := gitApi.Get()
				helper.ReportStatus(success, "Get Actions Public Key", true)
				if !success {
					os.Exit(1)
				}
				for _, secret := range Conf.Secrets {
					// Encrypt and "PUT" secret into remote repository
					epP := secret.Encrypt(&pubkey)
					gitApi := lib.GitApiFromRemote(&remote, epP, "")
					gitApi.EndpointReposSecrets()
					gitApi.In.Endpoint += "/" + secret.Name
					success := gitApi.Put()
					helper.ReportStatus(success, secret.Name, true)
				}
			}
		}
	},
}

func init() {
	repoSetCmd.AddCommand(repoSetSecretsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
