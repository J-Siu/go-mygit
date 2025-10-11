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
	"sync"

	"github.com/J-Siu/go-gitapi/v2"
	"github.com/J-Siu/go-gitapi/v2/repo"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

var repoSetWikiFalseCmd = &cobra.Command{
	Use:     "false " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"f"},
	Short:   "Set to false.",
	Long:    "Set to false.  " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		var info repo.Wiki // api
		info.Has = false
		// If no repo/dir specified in command line, add a ""
		if len(args) == 0 {
			args = []string{"."}
		}
		for _, workPath := range args {
			for _, remote := range global.Conf.MergedRemotes {
				wg.Add(1)
				var gitApi *gitapi.GitApi = remote.GetGitApi(&workPath, &info, global.Flag.Debug)
				gitApi.EndpointRepos()
				gitApi.SetPatch()
				if global.Flag.NoParallel {
					lib.RepoDo(gitApi, &wg, true, &global.Flag)
				} else {
					go lib.RepoDo(gitApi, &wg, true, &global.Flag)
				}
			}
		}
		wg.Wait()
	},
}

func init() {
	repoSetWikiCmd.AddCommand(repoSetWikiFalseCmd)
}
