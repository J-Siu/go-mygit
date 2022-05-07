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
	"sync"

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-mygit/lib"
	"github.com/spf13/cobra"
)

// Get repo public key
var repoGetPublickeyCmd = &cobra.Command{
	Use:     "publickey [repository ...]",
	Aliases: []string{"pk"},
	Short:   "get public key",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if len(args) == 0 {
			args = append(args, "")
		}
		for _, repo := range args {
			for _, remote := range Conf.MergedRemotes {
				// if remote.Vendor != gitapi.Vendor_Github {
				// 	fmt.Printf("%s:(%s) action secret not supported.\n", remote.Name, remote.Vendor)
				// } else {
				wg.Add(1)
				var info gitapi.RepoPublicKey
				gitApi := lib.GitApiFromRemote(&remote, &info, repo)
				gitApi.EndpointReposSecretsPubkey()
				go repoGetFunc(gitApi, &wg)
				// }
			}
			wg.Wait()
		}
	},
}

func init() {
	repoGetCmd.AddCommand(repoGetPublickeyCmd)
}
