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

	"github.com/J-Siu/go-gitcmd/v3/gitcmd"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/J-Siu/go-mygit/v3/lib"
	"github.com/spf13/cobra"
)

// Mass git push
var rootPushCmd = &cobra.Command{
	Use:     "push " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"p"}, // rootPushCmd
	Short:   "Git push",
	Long:    "Git push. " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out         = make(chan *string, 10)
			remoteQueue = new(lib.RemoteQueue)
			wg          sync.WaitGroup
		)
		if len(args) == 0 {
			args = []string{"."}
		}
		remoteQueue.New(&args, &global.Conf.MergedRemotes)
		go func() {
			for _, remote := range remoteQueue.Queue {
				var (
					gitCmd1  = new(gitcmd.GitCmd).New(*remote.WorkPath)
					gitCmd2  = new(gitcmd.GitCmd).New(*remote.WorkPath)
					options1 = []string{*remote.Name}
					options2 = options1
				)
				if global.Flag.PushAll {
					options1 = append(options1, "--all")
				}
				lib.GitRunWrapper(gitCmd1.Push(options1), *remote.WorkPath, &wg, global.Flag.NoParallel, global.Flag.NoTitle, out)
				if global.Flag.Tag {
					options2 = append(options2, "--tags")
					lib.GitRunWrapper(gitCmd2.Push(options2), *remote.WorkPath, &wg, global.Flag.NoParallel, global.Flag.NoTitle, out)
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
	rootCmd.AddCommand(rootPushCmd)
	rootPushCmd.Flags().BoolVarP(&global.Flag.PushAll, "all", "a", false, "Push all branches")
	rootPushCmd.Flags().BoolVarP(&global.Flag.Tag, "tags", "t", false, "Push with tags")
}
