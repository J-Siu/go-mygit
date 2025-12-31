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
	"os"
	"sync"

	"github.com/J-Siu/go-gitapi/v3/api"
	"github.com/J-Siu/go-gitapi/v3/base"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/J-Siu/go-mygit/v3/helper"
	"github.com/J-Siu/go-mygit/v3/lib"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var repoSetSecretCmd = &cobra.Command{
	Use:     "secret " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"s"},
	Short:   "Set action secret",
	Long:    "Set action secret. " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out = make(chan *string, 10)
			wg  sync.WaitGroup
		)
		// If no repo/dir specified in command line, add a ""
		if len(args) == 0 {
			args = []string{"."}
		}
		// --name/--value must be used together
		if (global.Flag.Secret.Name != "" && global.Flag.Secret.Value == "") || (global.Flag.Secret.Name == "" && global.Flag.Secret.Value != "") {
			ezlog.Err().M("-n/--name and -v/--value must be used together").Out()
			os.Exit(1)
		}
		go func() {
			for _, workPath := range args {
				for _, remote := range global.Conf.MergedRemotes {
					if remote.Vendor != base.VendorGithub {
						fmt.Printf("%s(%s) action secret not supported.\n", remote.Name, remote.Vendor)
					} else {
						// "GET" public key
						ezlog.Log().N(remote.Name).Out()
						var (
							property = remote.GitApiProperty(&workPath, global.Flag.Debug)
						)
						var secretsP *lib.Secrets
						if global.Flag.Secret.Name != "" && global.Flag.Secret.Value != "" {
							// Use command line value
							secretsP = &lib.Secrets{global.Flag.Secret}
						} else {
							// Use Conf secrets
							secretsP = &global.Conf.Secrets
						}
						// Use config secrets
						for _, secret := range *secretsP {
							ezlog.Log().N("secret").M(secret).Out()
							var (
								ga = new(api.EncryptedPair).New(property).Set(secret.Name, secret.Value)
							)
							helper.GitApiDoWrapper(ga, &global.Flag, &wg, out)
						}
					}
				}
			}
			wg.Wait()
			close(out)
		}()
		global.Flag.SingleLine = true
		global.Flag.StatusOnly = true
		for o := range out {
			ezlog.Log().Se().M(o).Out()
		}
	},
}

func init() {
	repoSetCmd.AddCommand(repoSetSecretCmd)
	repoSetSecretCmd.Flags().StringVarP(&global.Flag.Secret.Name, "name", "n", "", "Secret name")
	repoSetSecretCmd.Flags().StringVarP(&global.Flag.Secret.Value, "value", "v", "", "Secret value")
}
