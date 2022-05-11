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
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var rootPushCmd = &cobra.Command{
	Use:     "push " + lib.TXT_REPO_DIR_USE,
	Aliases: []string{"p"}, // rootPushCmd
	Short:   "Git push",
	Long:    "Git push. " + lib.TXT_REPO_DIR_LONG + lib.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if len(args) == 0 {
			args = []string{"."}
		}
		for _, workpath := range args {
			if helper.GitRoot(&workpath) == "" {
				helper.Report("is not a git repository.", workpath, true, true)
				return
			}

			// Create queue base on local remote
			var remoteLocal []string = *helper.GitRemote(&workpath, false)
			var remoteQueue []string
			for _, remote := range Conf.MergedRemotes {
				// Only add to queue if exist locally
				if helper.StrArrayPtrContain(&remoteLocal, &remote.Name) {
					remoteQueue = append(remoteQueue, remote.Name)
				}
			}

			for _, remote := range remoteQueue {
				// var fullpath string = *helper.FullPath(&workpath)
				// make a local copy of workpath for go routine
				var wp string = workpath
				options1 := []string{remote}
				if Flag.PushAll {
					wg.Add(1)
					options2 := append(options1, "--all")
					go lib.GitPush(&wp, &options2, &wg)
				} else {
					wg.Add(1)
					go lib.GitPush(&wp, &options1, &wg)
				}
				if Flag.PushTag {
					wg.Add(1)
					options2 := append(options1, "--tags")
					go lib.GitPush(&wp, &options2, &wg)
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
