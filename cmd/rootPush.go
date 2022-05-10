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

	"github.com/J-Siu/go-helper"
	"github.com/J-Siu/go-mygit/lib"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var rootPushCmd = &cobra.Command{
	Use:     "push " + lib.TXT_REPO_DIR_USE,
	Aliases: []string{"p"},
	Short:   "Git push",
	Long:    "Git push. " + lib.TXT_REPO_DIR_LONG + lib.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if len(args) == 0 {
			args = []string{*helper.CurrentPath()}
		}
		for _, workpath := range args {
			if helper.GitRoot(&workpath) == "" {
				helper.Report("is not a git repository.", workpath, true, true)
				return
			}
			for _, remote := range Conf.MergedRemotes {
				var fullpath string = *helper.FullPath(&workpath)
				args := []string{remote.Name}
				if Flag.PushAll {
					wg.Add(1)
					a := append(args, "--all")
					go lib.GitPush(&fullpath, &a, &wg)
				} else {
					wg.Add(1)
					go lib.GitPush(&fullpath, &args, &wg)
				}
				if Flag.PushTag {
					wg.Add(1)
					a := append(args, "--tags")
					go lib.GitPush(&fullpath, &a, &wg)
				}
			}
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(rootPushCmd)
	rootPushCmd.Flags().BoolVarP(&Flag.PushAll, "all", "a", false, "Push all branches")
	rootPushCmd.Flags().BoolVarP(&Flag.PushTag, "tags", "t", false, "Push all branches")
	// rootPushCmd.MarkFlagsMutuallyExclusive("all","tags")
}
