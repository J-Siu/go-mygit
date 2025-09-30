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
	"os"
	"sync"

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

// Mass git clone
var rootCloneCmd = &cobra.Command{
	Use:     "clone " + global.TXT_REPO_CLONE_USE,
	Aliases: []string{"cl"},
	Short:   "Git clone",
	Long:    "Git clone. " + global.TXT_REPO_CLONE_LONG,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		// Check flag
		if len(global.Flag.Groups) != 0 || // should not have --group
			len(global.Flag.Remotes) != 1 || // need exactly 1 --remote
			len(args) == 0 { // need >0 repo name
			ezlog.Log().M(global.TXT_REPO_CLONE_LONG).Out()
			os.Exit(1)
		}
		// Check remote name exist
		var remote *lib.Remote = global.Conf.Remotes.GetByName(&global.Flag.Remotes[0])
		if remote == nil {
			ezlog.Log().N("Remote not configured").M(global.Flag.Remotes[0]).Out()
			os.Exit(1)
		}

		for _, repoName := range args {
			wg.Add(1)

			// construct url
			var options []string = []string{remote.Ssh + ":/" + remote.User + "/" + repoName}
			if global.Flag.NoParallel {
				lib.GitClone(&options, &wg, global.Flag.NoTitle)
			} else {
				go lib.GitClone(&options, &wg, global.Flag.NoTitle)
			}
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(rootCloneCmd)
}
