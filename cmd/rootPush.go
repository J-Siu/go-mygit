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

	"github.com/J-Siu/go-gitcmd"

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

type remoteQueueStruct struct {
	Name     *string
	WorkPath *string
}

// Mass git push
var rootPushCmd = &cobra.Command{
	Use:     "push " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"p"}, // rootPushCmd
	Short:   "Git push",
	Long:    "Git push. " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out         = make(chan *string, 10)
			remoteLocal []string
			remoteQueue []*remoteQueueStruct
			wg          sync.WaitGroup
		)

		if len(args) == 0 {
			args = []string{"."}
		}
		for _, workPath := range args {
			// Create queue base on local remote
			remoteLocal = *gitcmd.GitRemote(&workPath, false)
			for _, remote := range global.Conf.MergedRemotes {
				// Only add to queue if exist locally
				if str.ArrayContains(&remoteLocal, &remote.Name, false) {
					rs := remoteQueueStruct{Name: &remote.Name, WorkPath: &workPath}
					remoteQueue = append(remoteQueue, &rs)
				} else {
					ezlog.Log().N(workPath).N("Remote not setup").M(remote.Name)
				}
			}
		}
		ezlog.Debug().N("RemoteQueue").Lm(remoteQueue).Out()

		go func() {
			for _, remote := range remoteQueue {
				var (
					options1 = []string{*remote.Name}
					options2 = options1
				)
				if global.Flag.PushAll {
					options2 = append(options1, "--all")
				}
				push(remote.WorkPath, &options2, &wg, out)
				if global.Flag.PushTag {
					options2 = append(options1, "--tags")
					push(remote.WorkPath, &options2, &wg, out)
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
	rootPushCmd.Flags().BoolVarP(&global.Flag.PushTag, "tags", "t", false, "Push with tags")
}

func push(wp *string, options *[]string, wg *sync.WaitGroup, out chan *string) {
	w := *wp
	opts := *options
	if global.Flag.NoParallel {
		lib.GitPush(&w, &opts, nil, global.Flag.NoTitle, out)
	} else {
		wg.Add(1)
		go lib.GitPush(&w, &opts, wg, global.Flag.NoTitle, out)
	}
}
