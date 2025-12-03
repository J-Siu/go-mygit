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

	"github.com/J-Siu/go-gitapi/v2/gitapi"
	"github.com/J-Siu/go-gitapi/v2/repo"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

// Get repo visibility
var repoGetActionsCmd = &cobra.Command{
	Use:     "actions " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"act", "action"},
	Short:   "Get wiki status",
	Long:    "Get wiki status. " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out = make(chan *string, 10)
			wg  sync.WaitGroup
		)
		if len(args) == 0 {
			args = []string{"."}
		}
		go func() {
			for _, workPath := range args {
				for _, remote := range global.Conf.MergedRemotes {
					var (
						info   gitapi.IInfo
						gitApi *gitapi.GitApi
					)
					if remote.Vendor == gitapi.VendorGithub {
						info = &repo.ActionsGithub{} // Github API
						gitApi = remote.GetGitApi(&workPath, info, global.Flag.Debug).EndpointReposActionsGithub()
					} else {
						info = &repo.Actions{} // Gitea API
						gitApi = remote.GetGitApi(&workPath, info, global.Flag.Debug).EndpointRepos()
					}

					gitApi.SetGet()
					lib.RepoDoRun(gitApi, global.Flag, true, false, &wg, out)
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
	repoGetCmd.AddCommand(repoGetActionsCmd)
}
