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

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-mygit/lib"
	"github.com/spf13/cobra"
)

// descriptionCmd represents the description command
var repoSetDescriptionCmd = &cobra.Command{
	Use:     "description",
	Aliases: []string{"d", "des", "desc"},
	Short:   "set description",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		var info gitapi.RepoDescription
		if len(args) > 0 {
			info.Description = args[0]
		}
		for _, remote := range Conf.MergedRemotes {
			wg.Add(1)
			gitApi := lib.GitApiFromRemote(&remote, &info, "")
			gitApi.EndpointRepos()
			go repoPatchFunc(gitApi, &wg)
		}
		wg.Wait()
	},
}

func init() {
	repoSetCmd.AddCommand(repoSetDescriptionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// descriptionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// descriptionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
