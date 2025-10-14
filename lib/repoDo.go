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

package lib

import (
	"errors"
	"sync"

	"github.com/J-Siu/go-gitapi/v2"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
)

type RepoDoProperty struct {
	GitApi     *gitapi.GitApi  `json:"GitApi"`
	NoSkip     bool            `json:"NoSkip"`
	NoTitle    bool            `json:"NoTitle"`
	NoParallel bool            `json:"NoParallel"`
	SingleLine bool            `json:"SingleLine"`
	StatusOnly bool            `json:"StatusOnly"` // Display api request status
	Wg         *sync.WaitGroup `json:"Wg"`
}

// parallel wrapper
func RepoDo(property *RepoDoProperty) {
	if property.NoParallel {
		repoDoProcess(property)
	} else {
		go repoDoProcess(property)
	}
}

func repoDoProcess(property *RepoDoProperty) {
	prefix := "RepoDo2"
	var (
		gitApi = property.GitApi
		wg     = property.Wg
		log    = new(ezlog.EzLog).New().SetLogLevel(ezlog.GetLogLevel())
	)
	if wg != nil {
		defer wg.Done()
	}
	var title string
	if !property.NoTitle {
		title = gitApi.Repo + "(" + gitApi.Name + ")"
	}

	status := gitApi.Do().Ok()
	if status {
		if property.StatusOnly {
			log.Log().N(title).Success(status).Out()
		} else {
			output := gitApi.Output()
			log.Debug().N(prefix).N("output").M(output).Out()
			if !(output == nil || *output == "") || property.NoSkip {
				if property.SingleLine {
					log.Log().N(title).M(output).Out()
				} else {
					log.Log().Nn(title).M(output).Out()
				}
			}
		}
	} else {
		// API or HTTP GET failed, try to extract error message
		log.Err().N(title).M(gitApi.Err())
		errs.Queue("", errors.New(log.String()))
	}
}
