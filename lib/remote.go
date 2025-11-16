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
	"path"
	"sync"

	"github.com/J-Siu/go-gitapi/v2"
	"github.com/J-Siu/go-gitcmd"
	"github.com/J-Siu/go-helper/v2/cmd"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
)

// Remote entry in config file
type Remote struct {
	Group  string        `json:"group,omitempty"`  // Group name
	Name   string        `json:"name,omitempty"`   // Name of remote entry, also use as git remote name
	Ssh    string        `json:"ssh,omitempty"`    // Ssh url for git server
	Vendor gitapi.Vendor `json:"vendor,omitempty"` // Api vendor/brand

	EntryPoint string `json:"entrypoint,omitempty"` // Api entrypoint url
	Private    bool   `json:"private,omitempty"`    // Default private value
	Token      string `json:"token,omitempty"`      // Api token.
	User       string `json:"user,omitempty"`       // Api user

	NoTitle    bool `json:"no_title,omitempty"`   // This is pass from global.Flag
	SkipVerify bool `json:"skipverify,omitempty"` // Api request skip cert verify (allow self-signed cert)
}

func (remote *Remote) GetGitApi(workPathP *string, info gitapi.IInfo, debug bool) *gitapi.GitApi {
	var fullPath string = *file.FullPath(workPathP)
	var repo string = path.Base(fullPath)

	property := gitapi.Property{
		Name:       remote.Name,
		Token:      remote.Token,
		EntryPoint: remote.EntryPoint,
		User:       remote.User,
		Vendor:     remote.Vendor,
		SkipVerify: remote.SkipVerify,
		Repo:       repo,
		Info:       info,
		Debug:      debug}

	apiP := gitapi.New(&property)
	// Set Github header
	apiP.HeaderGithub()
	return apiP
}

// Add all Remotes into git repository
func (remote *Remote) Add(workPathP *string) *cmd.Cmd {
	remote.Remove(workPathP)
	var fullPath string = *file.FullPath(workPathP)
	var repo string = path.Base(fullPath)
	var git string = remote.Ssh + ":/" + remote.User + "/" + repo + ".git"
	var myCmd = gitcmd.GitRemoteAdd(&fullPath, remote.Name, git)
	var title string
	if !remote.NoTitle {
		title = *workPathP + ": " + myCmd.CmdLn
	}
	myCmdLog(myCmd, title)
	return myCmd
}

// Remove all Remotes in git repository
func (remote *Remote) Remove(workPathP *string) *cmd.Cmd {
	var myCmd *cmd.Cmd
	if gitcmd.GitRemoteExist(workPathP, remote.Name) {
		myCmd = gitcmd.GitRemoteRemove(workPathP, remote.Name)
	}
	var title string
	if myCmd != nil && !remote.NoTitle {
		title = *workPathP + ": " + myCmd.CmdLn
		myCmdLog(myCmd, title)
	}
	return myCmd
}

// Push all Remotes in git repository
func GitPush(workPathP *string, optionsP *[]string, wgP *sync.WaitGroup, noTitle bool) *cmd.Cmd {
	if wgP != nil {
		defer wgP.Done()
	}
	var myCmd *cmd.Cmd = gitcmd.GitPush(workPathP, optionsP)
	var title string
	if !noTitle {
		title = *workPathP + ": " + myCmd.CmdLn
	}
	myCmdLog(myCmd, title)
	return myCmd
}

// Push all Remotes in git repository
func GitPull(workPathP *string, optionsP *[]string, wgP *sync.WaitGroup, noTitle bool) *cmd.Cmd {
	if wgP != nil {
		defer wgP.Done()
	}

	var myCmd *cmd.Cmd = gitcmd.GitPull(workPathP, optionsP)

	var title string
	if !noTitle {
		title = *workPathP + ": " + myCmd.CmdLn
	}
	myCmdLog(myCmd, title)
	return myCmd
}

// git clone to current directory
func GitClone(optionsP *[]string, wgP *sync.WaitGroup, noTitle bool) *cmd.Cmd {
	if wgP != nil {
		defer wgP.Done()
	}

	var myCmd *cmd.Cmd = gitcmd.GitClone(nil, optionsP)

	var title string
	if !noTitle {
		title = myCmd.CmdLn
	}
	myCmdLog(myCmd, title)
	return myCmd
}

func myCmdLog(myCmd *cmd.Cmd, title string) {
	se := ezlog.GetSkipEmpty()
	ezlog.SetSkipEmpty(false)
	if len(title) > 0 {
		ezlog.Log().N(title).Out()
	}
	ezlog.Log().Se().M(myCmd.Stdout.String()).Out()
	ezlog.Log().Se().M(myCmd.Stderr.String()).Out()
	ezlog.SetSkipEmpty(se)
}
