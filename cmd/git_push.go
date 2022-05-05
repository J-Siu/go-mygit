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
	"log"
	"os"
	"sync"

	"github.com/J-Siu/go-helper"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push to all remotes",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		if Flag.PushAll && Flag.PushTag {
			log.Fatal("--all/-a and --tags/-t cannot be used together")
			os.Exit(1)
		}
		for _, remote := range Conf.MergedRemotes {
			wg.Add(1)
			title := remote.Name
			args := []string{"push", title}
			if Flag.PushAll {
				args = append(args, "--all")
			}
			if Flag.PushTag {
				args = append(args, "--tags")
			}
			go helper.MyCmdRunWg("git", &args, &title, &wg, true)
		}
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(gitPushCmd)
	gitPushCmd.Flags().BoolVarP(&Flag.PushAll, "all", "a", false, "Push all branches")
	gitPushCmd.Flags().BoolVarP(&Flag.PushTag, "tags", "t", false, "Push all branches")
	// gitPushCmd.MarkFlagsMutuallyExclusive("all","tags")
}
