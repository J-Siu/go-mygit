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
	"os"

	"github.com/J-Siu/go-helper"
	"github.com/J-Siu/go-mygit/v2/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-mygit",
	Short: `Git and Repo automation made easy.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		helper.Debug = lib.Flag.Debug
		lib.Conf.Init(&lib.Flag)
		lib.Conf.MergeRemotes(&lib.Flag)
		if lib.Flag.Debug {
			helper.Report(&lib.Conf.File, "", true, true)
			helper.Report(&lib.Flag, "Flag", true, false)
			helper.Report(&lib.Conf.MergedRemotes, "Merged Remote", true, false)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.Debug, "debug", "d", false, "Enable debug")
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.NoSkip, "no-skip", "", false, "Don't skip empty output")
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.NoParallel, "no-parallel", "", false, "Don't process in parallel")
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.NoTitle, "no-title", "", false, "Don't print title for most output")
	rootCmd.PersistentFlags().StringArrayVarP(&lib.Flag.Groups, "group", "g", nil, "Specify group")
	rootCmd.PersistentFlags().StringArrayVarP(&lib.Flag.Remotes, "remote", "r", nil, "Specify remotes")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-mygit.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".go-mygit" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".go-mygit")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		helper.Report(err.Error(), "", true, true)
		os.Exit(1)
	}
}
