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
	"github.com/J-Siu/go-gitcmd"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var remoteListCmd = &cobra.Command{
	Use:     "list " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"l", "ls"},
	Short:   "List git remote",
	Long:    "List git remote. " + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = []string{"."}
		}
		for _, workPath := range args {
			if gitcmd.GitRoot(&workPath) == "" {
				ezlog.Log().N(workPath).M("is not a git repository").Out()
				return
			}
			var gitRemoteList *[]string = gitcmd.GitRemote(&workPath, true)
			var title string
			if !global.Flag.NoTitle {
				title = workPath
			}
			ezlog.Log().Nn(title).M(gitRemoteList).Out()
		}
	},
}

func init() {
	remoteCmd.AddCommand(remoteListCmd)
}
