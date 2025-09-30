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
	"strings"
	"sync"

	"github.com/J-Siu/go-gitcmd"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

// Mass git pull
var rootPullCmd = &cobra.Command{
	Use:   "pull " + global.TXT_REPO_DIR_USE,
	Short: "Git pull",
	Long:  "Git pull. " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if len(args) == 0 {
			args = []string{"."}
		}

		// Check flag
		if len(global.Flag.Groups) != 0 || // should not have --group
			len(global.Flag.Remotes) != 1 { // need exactly 1 --remote
			ezlog.Log().Msg(global.TXT_REPO_CLONE_LONG).Out()
			os.Exit(1)
		}
		// Check remote name exist
		var remote *lib.Remote = global.Conf.Remotes.GetByName(&global.Flag.Remotes[0])
		if remote == nil {
			ezlog.Log().Name("Remote not configured").Msg(global.Flag.Remotes[0]).Out()
			os.Exit(1)
		}

		for _, workPath := range args {
			if gitcmd.GitRoot(&workPath) == "" {
				ezlog.Log().Name(workPath).Msg("is not a git repository").Out()
				continue
			}

			var wp string = workPath
			var branch string = strings.TrimSpace(gitcmd.GitBranchCurrent(&wp).Stdout.String())
			var options []string = []string{remote.Name, branch}

			pull(&wp, &options, &wg)

			if global.Flag.PushTag {
				options = append(options, "--tags")
				pull(&wp, &options, &wg)
			}
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(rootPullCmd)
	// Re-use push flags
	rootPullCmd.Flags().BoolVarP(&global.Flag.PushTag, "tags", "t", false, "Pull all tags")
}

func pull(wp *string, options *[]string, wg *sync.WaitGroup) {
	wg.Add(1)
	w := *wp
	opts := *options
	if global.Flag.NoParallel {
		lib.GitPull(&w, &opts, wg, global.Flag.NoTitle)
	} else {
		go lib.GitPull(&w, &opts, wg, global.Flag.NoTitle)
	}
}
