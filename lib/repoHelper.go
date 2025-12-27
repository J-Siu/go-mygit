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

	"github.com/J-Siu/go-gitapi/v3/api"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
)

// `lib.RepoDoRun` wrapper
// Encapsulate RepoDo setup and run, sync.WaitGroup add amd done
func RepoDoRun(gitApi api.IApi, noParallel bool, wg *sync.WaitGroup, out chan api.IApi) {
	property := RepoDoProperty{
		GitApi:     gitApi,
		NoParallel: noParallel,
		Output:     out,
		Wg:         wg,
	}
	new(RepoDo).New(&property).Run()
}

func RepoOutput(gitApi api.IApi, flag TypeFlag, singleLine, statusOnly bool) {
	title := *gitApi.Repo() + "(" + gitApi.Name() + ")"
	status := gitApi.Ok()
	if status {
		if statusOnly {
			ezlog.Log()
			if !flag.NoTitle {
				ezlog.N(title)
			}
			ezlog.Success(status).Out()
		} else {
			output := gitApi.Output()
			if flag.NoSkip || (output != nil && *output != "") {
				ezlog.Log()
				if !flag.NoTitle {
					if singleLine {
						ezlog.N(title)
					} else {
						ezlog.Nl(title)
					}
				}
				ezlog.M(output).Out()
			}
		}
	} else {
		// API or HTTP GET failed, try to extract error message
		ezlog.Err().N(title).M(gitApi.Err())
		errs.Queue("", errors.New(ezlog.String()))
	}
}
