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
	"log"
	"os"
	"path"
	"sync"

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-helper"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var repoSetSecretCmd = &cobra.Command{
	Use:     "secret " + lib.TXT_REPO_DIR_USE,
	Aliases: []string{"s"},
	Short:   "Set action secret",
	Long:    "Set action secret. " + lib.TXT_REPO_DIR_LONG + lib.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		// If no repo/dir specified in command line, add a ""
		if len(args) == 0 {
			args = []string{"."}
		}
		// --name/--value must be used together
		if lib.Flag.Secret.Name == "" && lib.Flag.Secret.Value != "" ||
			lib.Flag.Secret.Name != "" && lib.Flag.Secret.Value == "" {
			log.Fatal("-n/--name and -v/--value must be used together")
			os.Exit(1)
		}
		for _, workPath := range args {
			for _, remote := range lib.Conf.MergedRemotes {
				if remote.Vendor != gitapi.Vendor_Github {
					fmt.Printf("%s(%s) action secret not supported.\n", remote.Name, remote.Vendor)
				} else {
					// "GET" public key
					helper.Report("", remote.Name, false, false)
					var pubkey gitapi.RepoPublicKey
					var gitApi *gitapi.GitApi = remote.GetGitApi(&workPath, &pubkey)
					gitApi.EndpointReposSecretsPubkey()
					ok := gitApi.SetGet().Do().Ok()
					helper.ReportStatus(ok, "Get Actions Public Key", true)
					if !ok {
						os.Exit(1)
					}
					// A list of secret to use
					var secretsP *lib.ConfSecrets
					if lib.Flag.Secret.Name != "" && lib.Flag.Secret.Value != "" {
						// Use command line value
						secretsP = &lib.ConfSecrets{lib.Flag.Secret}
					} else {
						// Use Conf secrets
						secretsP = &lib.Conf.Secrets
					}
					// Use config secrets
					for _, secret := range *secretsP {
						wg.Add(1)
						var infoP *gitapi.RepoEncryptedPair = secret.Encrypt(&pubkey)
						var gitApi *gitapi.GitApi = remote.GetGitApi(&workPath, infoP)
						gitApi.EndpointReposSecrets()
						gitApi.Req.Endpoint = path.Join(gitApi.Req.Endpoint, secret.Name)
						gitApi.SetPut()
						if lib.Flag.NoParallel {
							repoDo(gitApi, &wg, true)
						} else {
							go repoDo(gitApi, &wg, true)
						}
					}
				}
			}
		}
		wg.Wait()
	},
}

func init() {
	repoSetCmd.AddCommand(repoSetSecretCmd)
	repoSetSecretCmd.Flags().StringVarP(&lib.Flag.Secret.Name, "name", "n", "", "Secret name")
	repoSetSecretCmd.Flags().StringVarP(&lib.Flag.Secret.Value, "value", "v", "", "Secret value")
}
