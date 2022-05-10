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

func repoDelFunc(gitApi *gitapi.GitApi, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	gitApi.Del()
	helper.ReportStatus(gitApi.Res.Ok(), gitApi.Repo+"("+gitApi.Name+")", true)
}

func repoGetFunc(gitApi *gitapi.GitApi, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	// var success bool = gitApi.Get().Res.Ok()
	// if success {
	if gitApi.Get().Res.Ok() {
		// API GET OK
		var singleLine bool
		switch *gitApi.Res.Output {
		case "true", "false", "public", "private":
			singleLine = true
		default:
			singleLine = false
		}
		helper.Report(gitApi.Res.Output, gitApi.Repo+"("+gitApi.Name+")", !Flag.NoSkip, singleLine)
	} else {
		// API GET failed, try to extract error message
		var info errMsg
		err := json.Unmarshal([]byte(*gitApi.Res.Output), &info)
		if err == nil {
			helper.Report(info.Message, gitApi.Repo+"("+gitApi.Name+")", true, true)
		} else {
			helper.Report(gitApi.Res.Output, gitApi.Repo+"("+gitApi.Name+")", true, false)
		}
	}
}

func repoPatchFunc(gitApi *gitapi.GitApi, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	helper.ReportStatus(gitApi.Patch().Res.Ok(), gitApi.Repo+"("+gitApi.Name+")", true)
}

func repoPostFunc(gitApi *gitapi.GitApi, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	helper.ReportStatus(gitApi.Post().Res.Ok(), gitApi.Repo+"("+gitApi.Name+")", true)
}

func repoPutFunc(gitApi *gitapi.GitApi, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	helper.ReportStatus(gitApi.Put().Res.Ok(), gitApi.Repo+"("+gitApi.Name+")", true)
}
