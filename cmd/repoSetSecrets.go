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
	"log"
	"os"

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-helper"
	"github.com/J-Siu/go-mygit/lib"
	"github.com/spf13/cobra"
)

var name string
var value string

// setCmd represents the set command
var repoSetSecretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"s", "secret"},
	Short:   "set action secrets",
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" && value != "" || name != "" && value == "" {
			log.Fatal("-n/--name and -v/--value must be used together")
			os.Exit(1)
		}
		for _, remote := range Conf.MergedRemotes {
			if remote.Vendor != gitapi.Vendor_Github {
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

				if name != "" && value != "" {
					// Use command line value
					var secret lib.ConfSecret
					secret.Name = name
					secret.Value = value
					epP := secret.Encrypt(&pubkey)
					gitApi := lib.GitApiFromRemote(&remote, epP, "")
					repoPutSecret(gitApi, nil, &secret.Name)
				} else {
					// Use config secrets
					for _, secret := range Conf.Secrets {
						// Encrypt and "PUT" secret into remote repository
						epP := secret.Encrypt(&pubkey)
						gitApi := lib.GitApiFromRemote(&remote, epP, "")
						repoPutSecret(gitApi, nil, &secret.Name)
					}
				}
			}
		}
	},
}

func init() {
	repoSetCmd.AddCommand(repoSetSecretsCmd)
	repoSetSecretsCmd.Flags().StringVarP(&name, "name", "n", "", "Secret name")
	repoSetSecretsCmd.Flags().StringVarP(&value, "value", "v", "", "Secret value")
}
