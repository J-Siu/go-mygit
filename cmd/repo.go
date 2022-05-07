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
	"encoding/json"
	"sync"

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-helper"
	"github.com/spf13/cobra"
)

type errMsg struct {
	Errors  string `json:"errors"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"r", "rep", "repo"},
	Short:   "Repository commands",
}

func init() {
	rootCmd.AddCommand(repoCmd)
}

func repoDelFunc[T gitapi.GitApiInfo](gitApi *gitapi.GitApi[T], wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	success := gitApi.Del()
	helper.ReportStatus(success, gitApi.Name+"("+gitApi.Repo+")", true)
}

func repoGetFunc[T gitapi.GitApiInfo](gitApi *gitapi.GitApi[T], wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	var success bool = gitApi.Get()

	if success {
		var singleLine bool
		switch *gitApi.Out.Output {
		case "true", "false", "public", "private":
			singleLine = true
		default:
			singleLine = false
		}
		helper.Report(gitApi.Out.Output, gitApi.Name+"("+gitApi.Repo+")", true, singleLine)
	} else {
		// API failed, try to extract error message
		var info errMsg
		err := json.Unmarshal([]byte(*gitApi.Out.Output), &info)
		if err == nil {
			helper.Report(info.Message, gitApi.Name+"("+gitApi.Repo+")", true, true)
		} else {
			helper.Report(gitApi.Out.Output, gitApi.Name+"("+gitApi.Repo+")", true, false)
		}
	}
}

func repoPatchFunc[T gitapi.GitApiInfo](gitApi *gitapi.GitApi[T], wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	success := gitApi.Patch()
	helper.ReportStatus(success, gitApi.Name+"("+gitApi.Repo+")", true)
}

func repoPostFunc[T gitapi.GitApiInfo](gitApi *gitapi.GitApi[T], wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	success := gitApi.Post()
	helper.ReportStatus(success, gitApi.Name+"("+gitApi.Repo+")", true)
}

func repoPutFunc[T gitapi.GitApiInfo](gitApi *gitapi.GitApi[T], wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	success := gitApi.Put()
	helper.ReportStatus(success, gitApi.Name+"("+gitApi.Repo+")", true)
}

// func repoPutSecret[T gitapi.GitApiInfo](gitApi *gitapi.GitApi[T], wg *sync.WaitGroup) {
// 	if wg != nil {
// 		defer wg.Done()
// 	}
// 	gitApi.Put()
// 	success := gitApi.Put()
// 	helper.ReportStatus(success, gitApi.Name+"("+gitApi.Repo+")", true)
// }

func repoPutSecret[T gitapi.GitApiInfo](gitApi *gitapi.GitApi[T], wg *sync.WaitGroup, secretName *string) {
	if wg != nil {
		defer wg.Done()
	}
	gitApi.EndpointReposSecrets()
	gitApi.In.Endpoint += "/" + *secretName
	success := gitApi.Put()
	helper.ReportStatus(success, *secretName, true)
}
