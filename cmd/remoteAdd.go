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
	"github.com/J-Siu/go-gitcmd/v3/gitcmd"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/J-Siu/go-mygit/v3/helper"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var remoteAddCmd = &cobra.Command{
	Use:     "add " + global.TXT_REPO_DIR_USE,
	Aliases: []string{"a"},
	Short:   "Add git remote",
	Long:    "Add git remote. " + global.TXT_FLAGS_USE,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			out = make(chan *string, 10)
		)
		if len(args) == 0 {
			args = []string{"."}
		}
		global.Flag.NoParallel = true
		go func() {
			for _, workPath := range args {
				if gitcmd.Root(workPath) == "" {
					ezlog.Log().N(workPath).M("is not a git repository").Out()
					continue
				}
				for _, remote := range global.Conf.MergedRemotes {
					var (
						gitCmdRun1 = new(helper.GitCmdRun)
						property1  = new(helper.GitCmdRunProperty)
					)
					*property1 = helper.GitCmdRunProperty{
						Flag:     &global.Flag,
						OutChan:  out,
						Wg:       nil,
						WorkPath: workPath,
					}
					gitCmdRun1.New(property1).GitCmd.RemoteAdd(remote.Name, remote.GitUrl(workPath))
					gitCmdRun1.Run()
				}
			}
			close(out)
		}()
		for o := range out {
			ezlog.Log().Se().M(o).Out()
		}
	},
}

func init() {
	remoteCmd.AddCommand(remoteAddCmd)
}
