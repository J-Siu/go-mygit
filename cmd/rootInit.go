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

	"github.com/J-Siu/go-gitcmd/v3/gitcmd"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/spf13/cobra"
)

// Mass git init
var rootInitCmd = &cobra.Command{
	Use:     "init " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"i"}, // rootInitCmd
	Short:   "Git init and reset remotes",
	Long:    "Git init. Reset and add remotes. " + global.TXT_REPO_DIR_LONG + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out = make(chan *string, 10)
		)
		if len(args) == 0 {
			args = []string{"."}
		}
		go func() {
			for _, workPath := range args {
				gitcmd.Init(workPath)
				gitcmd.RemoteRemoveAll(workPath)
				for _, remote := range global.Conf.MergedRemotes {
					remote.Output = out
					remote.Add(&workPath)
				}
			}
			close(out)
		}()
		for o := range out {
			fmt.Print(*o)
		}
	},
}

func init() {
	rootCmd.AddCommand(rootInitCmd)
}
