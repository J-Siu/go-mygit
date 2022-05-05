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

package lib

import (
	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-helper"
)

type Remote struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	Ssh   string `json:"ssh"`

	Entrypoint string `json:"entrypoint"` // Api entrypoint url
	Token      string `json:"token"`      // Api token
	User       string `json:"user"`       // Api user
	Private    bool   `json:"private"`    // Default private vaule
	Vendor     string `json:"vendor"`     // Api vendor/brand
}

// Return a GitApi pointer base on Remote
func GitApiFromRemote[T gitapi.GitApiInfo](remoteP *Remote, info T) *gitapi.GitApi[T] {
	return gitapi.GitApiNew(
		remoteP.Name,
		remoteP.Token,
		remoteP.Entrypoint,
		remoteP.User,
		remoteP.Vendor,
		info)
}

// Add all Remotes into git repository
func (self *Remote) GitAdd() *helper.MyCmd {
	self.GitRemove()
	repo := helper.CurrentDirBase()
	git := self.Ssh + ":/" + self.User + "/" + repo + ".git"
	return helper.GitRemoteAdd(self.Name, git)
}

// Remove all Remotes in git repository
func (self *Remote) GitRemove() *helper.MyCmd {
	return helper.GitRemoteRemove(self.Name)
}
