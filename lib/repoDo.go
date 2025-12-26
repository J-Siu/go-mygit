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
	"sync"

	"github.com/J-Siu/go-gitapi/v3/api"
	"github.com/J-Siu/go-helper/v2/basestruct"
)

type RepoDoProperty struct {
	GitApi     api.IApi        `json:"GitApi"`
	NoParallel bool            `json:"NoParallel"`
	Output     chan api.IApi   `json:"Output"`
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

// Encapsulate wait group add amd done
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

func (t *RepoDo) process() *RepoDo {
	if t.Wg != nil {
		defer t.Wg.Done()
	}
	t.GitApi.Do()
	t.Output <- t.GitApi
	return t
}
