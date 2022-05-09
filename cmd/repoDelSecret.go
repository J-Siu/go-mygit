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
	"sync"

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-mygit/lib"
	"github.com/spf13/cobra"
)

// Delete repository action secret
var repoDelSecretCmd = &cobra.Command{
	Use:     "secret [repository ...]",
	Aliases: []string{"s"},
	Short:   "Delete remote repository secret",
	Long:    "Delete remote repository secret. If no repository is specified, current git root will be used as repository name. If --name, --value are not set, all secrets in conf will be added.",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		// Use secrets in conf if not specified in command line
		if len(Flag.SecretsDel) == 0 {
			for _, s := range Conf.Secrets {
				Flag.SecretsDel = append(Flag.SecretsDel, s.Name)
			}
		}

		// If no repo specified in command line, add a ""
		if len(args) == 0 {
			args = []string{""}
		}

		for _, repo := range args {
			for _, remote := range Conf.MergedRemotes {
				for _, secret := range Flag.SecretsDel {
					if remote.Vendor != gitapi.Vendor_Github {
						fmt.Printf("%s: only (%s) action secret is supported.\n", remote.Name, gitapi.Vendor_Github)
					} else {
						wg.Add(1)
						gitApi := lib.GitApiFromRemote(&remote, gitapi.Nil(), repo)
						gitApi.EndpointReposSecrets()
						gitApi.In.Endpoint += "/" + secret
						go repoDelFunc(gitApi, &wg)
					}
				}
			}
		}
		wg.Wait()
	},
}

func init() {
	repoDelCmd.AddCommand(repoDelSecretCmd)
	repoDelSecretCmd.Flags().StringArrayVarP(&Flag.SecretsDel, "name", "n", []string{}, "Secret name")
}
