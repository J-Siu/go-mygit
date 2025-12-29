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
	"github.com/J-Siu/go-gitcmd/v3/gitcmd"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/str"
)

type RemoteQueueItem struct {
	Name     *string
	WorkPath *string
}

type RemoteQueue struct {
	Queue []*RemoteQueueItem
}

func (t *RemoteQueue) New(workPaths *[]string, remotes *Remotes) *RemoteQueue {
	var (
		remoteLocal []string
	)
	t.Queue = nil
	for _, workPath := range *workPaths {
		// Create queue base on local remote
		remoteLocal = *gitcmd.Remote(workPath, false)
		for _, remote := range *remotes {
			// Only add to queue if exist locally
			if str.ArrayContains(&remoteLocal, &remote.Name, false) {
				item := RemoteQueueItem{Name: &remote.Name, WorkPath: &workPath}
				t.Queue = append(t.Queue, &item)
			} else {
				ezlog.Log().N(workPath).N("Remote not setup").M(remote.Name)
			}
		}
	}
	ezlog.Debug().N("RemoteQueue").Lm(t).Out()
	return t
}
