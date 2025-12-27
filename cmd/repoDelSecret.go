/*
Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

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

	"github.com/J-Siu/go-gitapi/v3/api"
	"github.com/J-Siu/go-gitapi/v3/base"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/J-Siu/go-mygit/v3/lib"
	"github.com/spf13/cobra"
)

// Delete repository action secret
var repoDelSecretCmd = &cobra.Command{
	Use:     "secret " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"s"},
	Short:   "Delete action secret",
	Long:    "Delete action secret. If --name is not set, all secrets in config will be used. " + global.TXT_REPO_DIR_LONG,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out = make(chan api.IApi, 10)
			wg  sync.WaitGroup
		)
		// If no repo specified in command line, add a ""
		if len(args) == 0 {
			args = []string{"."}
		}
		// Use secrets in conf if not specified in command line
		if len(global.Flag.SecretsDel) == 0 {
			for _, s := range global.Conf.Secrets {
				global.Flag.SecretsDel = append(global.Flag.SecretsDel, s.Name)
			}
		}
		go func() {
			for _, workPath := range args {
				for _, remote := range global.Conf.MergedRemotes {
					if remote.Vendor != base.VendorGithub {
						fmt.Printf("%s(%s) action secret not supported.\n", remote.Name, remote.Vendor)
					} else {
						for _, secret := range global.Flag.SecretsDel {
							var (
								property = remote.GitApiProperty(&workPath, global.Flag.Debug)
								ga       = new(api.Repo).New(property).DelSecret(secret)
							)
							lib.RepoDoRun(ga, global.Flag.NoParallel, &wg, out)
						}
					}
				}
			}
			wg.Wait()
			close(out)
		}()
		for gitApi := range out {
			lib.RepoOutput(gitApi, global.Flag, true, true)
		}

	},
}

func init() {
	repoDelCmd.AddCommand(repoDelSecretCmd)
	repoDelSecretCmd.Flags().StringArrayVarP(&global.Flag.SecretsDel, "name", "n", []string{}, "Secret name")
}
