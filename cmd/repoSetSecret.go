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
	"path"
	"sync"

	"github.com/J-Siu/go-gitapi/v2/gitapi"
	"github.com/J-Siu/go-gitapi/v2/repo"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/J-Siu/go-mygit/v2/lib"

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
		if global.Flag.Secret.Name == "" || global.Flag.Secret.Value == "" {
			ezlog.Err().M("Both -n/--name and -v/--value must be set").Out()
			os.Exit(1)
		}
		go func() {
			for _, workPath := range args {
				for _, remote := range global.Conf.MergedRemotes {
					if remote.Vendor != gitapi.VendorGithub {
						fmt.Printf("%s(%s) action secret not supported.\n", remote.Name, remote.Vendor)
					} else {
						// "GET" public key
						ezlog.Log().N(remote.Name).Out()
						var pubkey repo.PublicKey
						var gitApi *gitapi.GitApi = remote.GetGitApi(&workPath, &pubkey, global.Flag.Debug)
						gitApi.EndpointReposSecretsPubkey()
						ok := gitApi.SetGet().Do().Ok()
						ezlog.Log().N("Get Actions Public Key").M(ok)
						if !ok {
							os.Exit(1)
						}
						// A list of secret to use
						var secretsP *lib.ConfSecrets
						if global.Flag.Secret.Name != "" && global.Flag.Secret.Value != "" {
							// Use command line value
							secretsP = &lib.ConfSecrets{global.Flag.Secret}
						} else {
							// Use Conf secrets
							secretsP = &global.Conf.Secrets
						}
						// Use config secrets
						for _, secret := range *secretsP {
							var infoP *repo.EncryptedPair = secret.Encrypt(&pubkey)
							var gitApi *gitapi.GitApi = remote.GetGitApi(&workPath, infoP, global.Flag.Debug)
							gitApi.EndpointReposSecrets()
							gitApi.Req.Endpoint = path.Join(gitApi.Req.Endpoint, secret.Name)
							gitApi.SetPut()
							lib.RepoDoRun(gitApi, global.Flag, true, true, &wg, out)
						}
					}
				}
			}
			wg.Wait()
			close(out)
		}()
		for o := range out {
			fmt.Print(*o)
		}

	},
}

func init() {
	repoSetCmd.AddCommand(repoSetSecretCmd)
	repoSetSecretCmd.Flags().StringVarP(&global.Flag.Secret.Name, "name", "n", "", "Secret name")
	repoSetSecretCmd.Flags().StringVarP(&global.Flag.Secret.Value, "value", "v", "", "Secret value")
}
