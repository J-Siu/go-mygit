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
	"strings"
	"sync"

	"github.com/J-Siu/go-gitcmd/v3/gitcmd"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/J-Siu/go-mygit/v3/lib"
	"github.com/spf13/cobra"
)

// Mass git pull
var rootPullCmd = &cobra.Command{
	Use:   "pull " + global.TXT_REPO_DIR_USE,
	Short: "Git pull",
	Long:  "Git pull. " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Check flag
		if len(global.Flag.Groups) != 0 || // should not have --group
			len(global.Flag.Remotes) != 1 { // need exactly 1 --remote
			ezlog.Log().M(global.TXT_REPO_CLONE_LONG).Out()
			os.Exit(1)
		}
		// Check remote name exist
		var remote *lib.Remote = global.Conf.Remotes.GetByName(&global.Flag.Remotes[0])
		if remote == nil {
			ezlog.Log().N("Remote not configured").M(global.Flag.Remotes[0]).Out()
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out    = make(chan *string, 10)
			remote = global.Conf.Remotes.GetByName(&global.Flag.Remotes[0])
			wg     sync.WaitGroup
		)
		if len(args) == 0 {
			args = []string{"."}
		}
		go func() {
			for _, workPath := range args {
				var (
					wp       = workPath
					branch   = strings.TrimSpace(gitcmd.BranchCurrent(wp).Stdout.String())
					gitCmd   = new(gitcmd.GitCmd).New(wp)
					options1 = []string{remote.Name, branch}
					options2 = options1
				)
				lib.GitRunWrapper(gitCmd.Pull(options1), wp, &wg, global.Flag.NoParallel, global.Flag.NoTitle, out)
				if global.Flag.Tag {
					options2 = append(options2, "--tags")
					lib.GitRunWrapper(gitCmd.Pull(options2), wp, &wg, global.Flag.NoParallel, global.Flag.NoTitle, out)
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
	rootCmd.AddCommand(rootPullCmd)
	// Re-use push flags
	rootPullCmd.Flags().BoolVarP(&global.Flag.Tag, "tags", "t", false, "Pull all tags")
}
