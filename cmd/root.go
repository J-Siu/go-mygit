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
	"os"

	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-mygit",
	Short:   `Git and Repo automation made easy.`,
	Version: global.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ezlog.SetLogLevel(ezlog.ErrLevel)
		if global.Flag.Debug {
			ezlog.SetLogLevel(ezlog.DebugLevel)
		}

		ezlog.Debug().N("Version").Mn(global.Version).Nn("Flag").M(&global.Flag).Out()

		global.Conf.New(&global.Flag.Groups, &global.Flag.Remotes)

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if !errs.IsEmpty() {
			ezlog.Err().Nn("Error").M(errs.Errs).Out()
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.Debug, "debug", "d", false, "Enable debug")
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.NoParallel, "no-parallel", "", false, "Don't process in parallel")
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.NoSkip, "no-skip", "", false, "Don't skip empty output")
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.NoTitle, "no-title", "", false, "Don't print title for most output")
	rootCmd.PersistentFlags().StringArrayVarP(&global.Flag.Groups, "group", "g", nil, "Specify group")
	rootCmd.PersistentFlags().StringArrayVarP(&global.Flag.Remotes, "remote", "r", nil, "Specify remotes")
	rootCmd.PersistentFlags().StringVarP(&global.Conf.FileConf, "config", "", lib.Default.FileConf, "Config file")
}
