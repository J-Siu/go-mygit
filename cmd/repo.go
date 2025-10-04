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
	"strconv"
	"sync"

	"github.com/J-Siu/go-gitapi/v2"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v2/global"
	"github.com/savaki/jq"
	"github.com/spf13/cobra"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"r", "rep", "repo"},
	Short:   "Repository commands",
}

func init() {
	rootCmd.AddCommand(repoCmd)
}

func repoDo(gitApi *gitapi.GitApi, wg *sync.WaitGroup, statusOnly bool) {
	if wg != nil {
		defer wg.Done()
	}
	var title string
	if !global.Flag.NoTitle {
		title = gitApi.Repo + "(" + gitApi.Name + ")"
	}

	status := gitApi.Do().Ok()
	if status {
		if statusOnly {
			// helper.ReportStatus(status, title, true)
			ezlog.Log().N(title).M(status).Out()
		} else {
			singleLine := false
			output := gitApi.Output()
			switch *output {
			case "true", "false", "public", "private":
				singleLine = true
			default:
				singleLine = false
			}
			if !(output == nil || *output == "") || global.Flag.NoSkip {
				if singleLine {
					ezlog.Log().N(title).M(output).Out()
				} else {
					ezlog.Log().Nn(title).M(output).Out()
				}
			}
		}
	} else {
		// API or HTTP GET failed, try to extract error message
		ezlog.Err().N(title).M(gitApi.Err()).Out()
	}
}

func repoUnarchiveGithub(gitApi *gitapi.GitApi, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	var title string
	if !global.Flag.NoTitle {
		title = gitApi.Repo + "(" + gitApi.Name + ")"
	}

	// Get Repo Node_Id and Name
	var info RepoNodeId
	gitApi.EndpointRepos()
	gitApi.Info = &info
	// Cannot use repoGetFunc(gitApi, nil), as we don't want print out here
	gitApi.Get()
	if !gitApi.Ok() {
		ezlog.Err().N(title).M(gitApi.Err()).Out()
		return
	}
	ezlog.Debug().Nn("RepoNodeId").M(&info).Out()

	// Use Github GraphQL as unarchive not supported by rest api
	gitApi.Req.EntryPoint = "https://api.github.com/graphql" // Github GraphQL entrypoint
	gitApi.Req.Endpoint = ""                                 // No endpoint for GraphGL
	gitApi.Info = nil                                        // Not using struct, not a REST operation
	// GraphQL query
	gitApi.Req.Data = `{
		"query":"mutation UnArchiveRepository($mutationId:String!,$repoID:ID!){unarchiveRepository(input:{clientMutationId:$mutationId,repositoryId:$repoID}){repository{isArchived}}}",
		"variables":{
			"mutationId":"true",
			"repoID":"` + info.Node_Id + `"
		}
	}`
	// repoPostFunc(gitApi, nil)
	gitApi.Post()

	var op jq.Op
	var err error

	// Response: {"data":{"unarchiveRepository":{"repository":{"isArchived":false}}}}
	// Use jq to extract isArchived and description
	op, err = jq.Parse(".data.unarchiveRepository.repository.isArchived")
	if err != nil {
		ezlog.Err().N(title).M(err).Out()
		return
	}
	isArchived, err := op.Apply(*gitApi.Res.Body)
	if err != nil {
		ezlog.Err().N(title).M(err).Out()
		return
	}
	// No error, print result
	ezlog.Log().N(title).M(isArchived).Out()
}

type RepoArchived struct {
	Archived bool `json:"archived"`
}

func (s *RepoArchived) String() string {
	return strconv.FormatBool(s.Archived)
}

func (s *RepoArchived) StringP() *string {
	tmp := s.String()
	return &tmp
}

type RepoNodeId struct {
	Name    string `json:"name"`
	Node_Id string `json:"node_id"`
}

func (s *RepoNodeId) StringP() *string {
	return &s.Node_Id
}

func (s *RepoNodeId) String() string {
	return *s.StringP()
}
