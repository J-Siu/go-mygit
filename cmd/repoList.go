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

	"github.com/J-Siu/go-gitapi/v3/api"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/J-Siu/go-mygit/v3/helper"
	"github.com/spf13/cobra"
)

// repo list
var repoListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all remote repositories",
	Long:    "List all remote repositories. " + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			// out = make(chan *string, 10)
			out = make(chan *string, 10)
			wg  sync.WaitGroup
		)
		go func() {
			for _, remote := range global.Conf.MergedRemotes {
				for p := 1; p < global.Flag.Page+1; p++ {
					var (
						property = remote.GitApiProperty(nil, global.Flag.Debug)
						ga       = new(api.InfoList).New(property, remote.Vendor, p).Get()
					)
					helper.GitApiDoWrapper(ga, &global.Flag, &wg, out)
				}
			}
			wg.Wait()
			close(out)
		}()
		global.Flag.SingleLine = false
		global.Flag.StatusOnly = false
		for o := range out {
			ezlog.Log().Se().M(o).Out()
		}
	},
}

func init() {
	repoCmd.AddCommand(repoListCmd)
	repoListCmd.Flags().IntVarP(&global.Flag.Page, "page", "p", 1, "Page number")
}
