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

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/J-Siu/go-mygit/v3/helper"
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
					gitCmdRun1 = new(helper.GitCmdRun)
					gitCmdRun2 = new(helper.GitCmdRun)
					options1   = []string{*remote.Name}
					options2   = append(options1, "--tags")
					property1  = new(helper.GitCmdRunProperty)
					property2  = new(helper.GitCmdRunProperty)
				)
				*property1 = helper.GitCmdRunProperty{
					Flag:     &global.Flag,
					OutChan:   out,
					Wg:       &wg,
					WorkPath: *remote.WorkPath,
				}
				if global.Flag.PushAll {
					options1 = append(options1, "--all")
				}
				gitCmdRun1.New(property1).GitCmd.Push(options1)
				gitCmdRun1.RunWrapper()

				if global.Flag.Tag {
					*property2 = *property1
					gitCmdRun2.New(property2).GitCmd.Push(options2)
					gitCmdRun2.RunWrapper()
				}
			}
			wg.Wait()
			close(out)
		}()
		for o := range out {
			ezlog.Log().Se().M(o).Out()
		}
	},
}

func init() {
	rootCmd.AddCommand(rootPushCmd)
	rootPushCmd.Flags().BoolVarP(&global.Flag.PushAll, "all", "a", false, "Push all branches")
	rootPushCmd.Flags().BoolVarP(&global.Flag.Tag, "tags", "t", false, "Push with tags")
}
