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

	"github.com/J-Siu/go-gitapi/v2/gitapi"
	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
)

type RepoDoProperty struct {
	GitApi     *gitapi.GitApi  `json:"GitApi"`
	NoParallel bool            `json:"NoParallel"`
	NoSkip     bool            `json:"NoSkip"`
	NoTitle    bool            `json:"NoTitle"`
	Output     chan *string    `json:"Output"`
	SingleLine bool            `json:"SingleLine"`
	StatusOnly bool            `json:"StatusOnly"` // Display api request status
	Wg         *sync.WaitGroup `json:"Wg"`
}

type RepoDo struct {
	*basestruct.Base
	*RepoDoProperty
}

func (t *RepoDo) New(property *RepoDoProperty) *RepoDo {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "Repo"

	t.RepoDoProperty = property
	return t
}

func (t *RepoDo) Run() *RepoDo {
	if t.NoParallel {
		t.Wg = nil
		t.process()
	} else {
		t.Wg.Add(1)
		go t.process()
	}
	return t
}

func (t *RepoDo) process() {
	prefix := t.MyType + ".process"
	var (
		gitApi = t.GitApi
		wg     = t.Wg
		log    = new(ezlog.EzLog).New().SetLogLevel(ezlog.GetLogLevel())
	)
	if wg != nil {
		defer wg.Done()
	}

	title := gitApi.Repo + "(" + gitApi.Name + ")"
	status := gitApi.Do().Ok()
	if status {
		if t.StatusOnly {
			log.Log()
			if !t.NoTitle {
				log.N(title)
			}
			t.Output <- log.Success(status).L().StringP()
		} else {
			output := gitApi.Output()
			t.Output <- log.Debug().N(prefix).N("output").M(output).StringP()
			if !(output == nil || *output == "") || t.NoSkip {
				log.Log()
				if !t.NoTitle {
					if t.SingleLine {
						log.N(title)
					} else {
						log.Nl(title)
					}
				}
				t.Output <- log.M(output).L().StringP()
			}
		}
	} else {
		// API or HTTP GET failed, try to extract error message
		log.Err().N(title).M(gitApi.Err())
		errs.Queue("", errors.New(log.String()))
	}
}

// `lib.RepoDoRun` wrapper
func RepoDoRun(gitApi *gitapi.GitApi, flag TypeFlag, singleLine, statusOnly bool, wg *sync.WaitGroup, out chan *string) {
	property := RepoDoProperty{
		GitApi:     gitApi,
		NoParallel: flag.NoParallel,
		NoSkip:     flag.NoSkip,
		NoTitle:    flag.NoTitle,
		Output:     out,
		SingleLine: singleLine,
		StatusOnly: statusOnly,
		Wg:         wg,
	}
	new(RepoDo).New(&property).Run()
}
