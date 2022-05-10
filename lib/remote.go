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
	"path"
	"sync"

	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-helper"
)

// Remote entry in config file
type Remote struct {
	Name  string `json:"name"`  // Name of remote entry, also use as git remote name
	Group string `json:"group"` // Group name
	Ssh   string `json:"ssh"`   // Ssh url for git server

	Entrypoint string `json:"entrypoint"` // Api entrypoint url
	Token      string `json:"token"`      // Api token
	User       string `json:"user"`       // Api user
	Private    bool   `json:"private"`    // Default private vaule
	Vendor     string `json:"vendor"`     // Api vendor/brand
}

func (self *Remote) GetGitApi(workpathP *string, info gitapi.GitApiInfo) *gitapi.GitApi {
	var fullpath string = *helper.FullPath(workpathP)
	var repo string = path.Base(fullpath)
	apiP := gitapi.GitApiNew(
		self.Name,
		self.Token,
		self.Entrypoint,
		self.User,
		self.Vendor,
		repo,
		info)
	// Set Github header
	apiP.HeaderGithub()
	return apiP
}

// Add all Remotes into git repository
func (self *Remote) GitAdd(workpathP *string) *helper.MyCmd {
	self.GitRemove(workpathP)
	var fullpath string = *helper.FullPath(workpathP)
	var gitroot string = helper.GitRoot(workpathP)
	if gitroot == "" {
		helper.Report("is not a git repository.", *workpathP, true, true)
		return nil
	}
	var repo string = path.Base(gitroot)
	var git string = self.Ssh + ":/" + self.User + "/" + repo + ".git"
	var myCmd *helper.MyCmd = helper.GitRemoteAdd(&fullpath, self.Name, git)
	var title string = *workpathP + ": " + myCmd.CmdLn
	helper.Report(myCmd.Stderr.String(), title, true, false)
	helper.Report(myCmd.Stdout.String(), title, true, false)
	return myCmd
}

// Remove all Remotes in git repository
func (self *Remote) GitRemove(workpathP *string) *helper.MyCmd {
	var myCmd *helper.MyCmd = helper.GitRemoteRemove(workpathP, self.Name)
	var title string = *workpathP + ": " + myCmd.CmdLn
	helper.Report(myCmd.Stderr.String(), title, true, false)
	helper.Report(myCmd.Stdout.String(), title, true, false)
	return myCmd
}

// Push all Remotes in git repository
func GitPush(workpathP *string, optionsP *[]string, wgP *sync.WaitGroup) *helper.MyCmd {
	if wgP != nil {
		defer wgP.Done()
	}
	var myCmd *helper.MyCmd = helper.GitPush(workpathP, optionsP)
	var title string = *workpathP + ": " + myCmd.CmdLn
	helper.Report(myCmd.Stderr.String(), title, true, false)
	helper.Report(myCmd.Stdout.String(), title, true, false)
	return myCmd
}
