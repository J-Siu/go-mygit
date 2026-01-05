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
	"sync"

	"github.com/J-Siu/go-gitcmd/v3/gitcmd"
	"github.com/J-Siu/go-helper/v2/cmd"
	"github.com/J-Siu/go-mygit/v3/lib"
)

type GitCmdRunProperty struct {
	Flag     *lib.TypeFlag `json:"Flag"`
	OutChan  chan *string
	Wg       *sync.WaitGroup
	WorkPath string
}

type GitCmdRun struct {
	*GitCmdRunProperty
	*gitcmd.GitCmd
}

func (t *GitCmdRun) New(property *GitCmdRunProperty) *GitCmdRun {
	t.GitCmdRunProperty = property
	t.GitCmd = new(gitcmd.GitCmd).New(t.WorkPath)
	return t
}

// handle goroutines and output
func (t *GitCmdRun) RunWrapper() *GitCmdRun {
	if t.Flag.NoParallel || t.Wg == nil {
		t.Wg = nil
		t.Run()
	} else {
		t.Wg.Add(1)
		go t.Run()
	}
	return t
}

func (t *GitCmdRun) Run() *cmd.Cmd {
	if t.Wg != nil {
		defer t.Wg.Done()
	}
	var myCmd = t.GitCmd.Run()
	t.OutChan <- t.Out()
	return myCmd
}

func (t *GitCmdRun) Out() *string {
	var (
		out string
		tmp string
	)
	if !t.Flag.NoTitle {
		out = t.GitCmd.CmdLn
		if t.WorkPath != "" {
			out = t.WorkPath + ": " + t.GitCmd.CmdLn
		}
	}
	if len(out) > 0 {
		out = out + ":\n"
	}
	tmp = t.GitCmd.Stdout.String()
	if tmp != "" {
		out += tmp
	}
	tmp = t.GitCmd.Stderr.String()
	if tmp != "" {
		out += tmp
	}
	return &out
}
