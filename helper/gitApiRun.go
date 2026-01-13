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

package helper

import (
	"errors"
	"sync"

	"github.com/J-Siu/go-gitapi/v3/api"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-mygit/v3/lib"
)

type GitApiRunProperty struct {
	Flag    *lib.TypeFlag   `json:"Flag"`
	IApi    api.IApi        `json:"IApi"`
	OutChan chan *string    `json:"Output"`
	Wg      *sync.WaitGroup `json:"Wg"`
}

type GitApiRun struct {
	*GitApiRunProperty
}

func (t *GitApiRun) New(property *GitApiRunProperty) *GitApiRun {
	t.GitApiRunProperty = property
	return t
}

// handle goroutines and output
func (t *GitApiRun) Run() *GitApiRun {
	if t.Flag.NoParallel || t.Wg == nil {
		t.run()
	} else {
		t.Wg.Go(t.run)
	}
	return t
}

func (t *GitApiRun) run() {
	t.IApi.Do()
	t.OutChan <- t.Out()
}

func (t *GitApiRun) Out() *string { //flag lib.TypeFlag, singleLine, statusOnly bool
	var (
		log    = new(ezlog.EzLog).New()
		status = t.IApi.Ok()
		title  = *t.IApi.Repo() + "(" + t.IApi.Name() + ")"
	)
	if status {
		log.Log()
		if t.Flag.StatusOnly {
			if !t.Flag.NoTitle {
				log.N(title)
			}
			log.Success(status)
		} else {
			output := t.IApi.Output()
			if t.Flag.NoSkip || (output != nil && *output != "") {
				if !t.Flag.NoTitle {
					if t.Flag.SingleLine {
						log.N(title)
					} else {
						log.Nl(title)
					}
				}
				log.M(output)
			}
		}
		return log.StringP()
	} else {
		// API or HTTP GET failed, try to extract error message
		ezlog.Err().N(title).M(t.IApi.Err())
		errs.Queue("", errors.New(ezlog.String()))
		return nil
	}
}

// Encapsulate GitApiDo setup and run, sync.WaitGroup add amd done
func GitApiRunWrapper(flag *lib.TypeFlag, wg *sync.WaitGroup, out chan *string, gitApi api.IApi) {
	property := GitApiRunProperty{
		Flag:    flag,
		IApi:    gitApi,
		OutChan: out,
		Wg:      wg,
	}
	new(GitApiRun).New(&property).Run()
}
