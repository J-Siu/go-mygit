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

	"github.com/J-Siu/go-gitapi/v3/base"
	"github.com/J-Siu/go-helper/v2/file"
)

// Remote entry in config file
type Remote struct {
	Group  string      `json:"group,omitempty"`  // Group name
	Name   string      `json:"name,omitempty"`   // Name of remote entry, also use as git remote name
	Ssh    string      `json:"ssh,omitempty"`    // Ssh url for git server
	Vendor base.Vendor `json:"vendor,omitempty"` // Api vendor/brand

	EntryPoint string `json:"entrypoint,omitempty"` // Api entrypoint url
	Private    bool   `json:"private,omitempty"`    // Default private value
	Token      string `json:"token,omitempty"`      // Api token.
	User       string `json:"user,omitempty"`       // Api user

	NoTitle    bool `json:"no_title,omitempty"`   // This is pass from global.Flag
	SkipVerify bool `json:"skipverify,omitempty"` // Api request skip cert verify (allow self-signed cert)

	Output chan *string `json:"-"` // Must mark json ignore
}

func (t *Remote) GitApiProperty(workPathP *string, debug bool) *base.Property {
	var (
		fullPath string = *file.FullPath(workPathP)
		repo     string = path.Base(fullPath)
	)
	property := base.Property{
		Name:       t.Name,
		Token:      t.Token,
		EntryPoint: t.EntryPoint,
		User:       t.User,
		Vendor:     t.Vendor,
		SkipVerify: t.SkipVerify,
		Repo:       repo,
		Debug:      debug,
	}
	return &property
}
