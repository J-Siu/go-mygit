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
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v3/global"
	"github.com/spf13/cobra"
)

// remoteCmd represents the remote command
var confAllCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "Print all configurations",
	Run: func(cmd *cobra.Command, args []string) {
		if global.FlagConf.Secret {
			ezlog.Log().N("Config").Lm(global.Conf.Remotes).Out()
		} else {
			ezlog.Log().N("Config").Lm(confSafe()).Out()
		}
	},
}

func init() {
	configCmd.AddCommand(confAllCmd)
}
