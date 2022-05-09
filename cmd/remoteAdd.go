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
	"github.com/J-Siu/go-helper"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var remoteAddCmd = &cobra.Command{
	Use:     "add [repository ...]",
	Aliases: []string{"a"},
	Short:   "Add git remotes base on configuration and flags",
	Long:    "Add git remotes base on configuration and flags. If -r/-g not specified, all remotes in config file will be added.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = []string{*helper.CurrentPath()}
		}
		for _, path := range args {
			if helper.GitRoot(&path) == "" {
				helper.Report("is not a git repository.", path, true, true)
				return
			}
			for _, remote := range Conf.MergedRemotes {
				remote.GitAdd(&path)
			}
		}
	},
}

func init() {
	remoteCmd.AddCommand(remoteAddCmd)
}
