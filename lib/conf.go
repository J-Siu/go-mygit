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
	"log"
	"os"

	"github.com/J-Siu/go-helper"
	"github.com/spf13/viper"
)

/*
File, Group, MergeRemotes are filled in at runtime

Remotes, Secrets are read from config file by viper
*/
type TypeConf struct {
	File          string      `json:"-"`
	Groups        Groups      `json:"-"`
	Remotes       Remotes     `json:"remotes"`
	Secrets       ConfSecrets `json:"secrets"`
	MergedRemotes Remotes     `json:"-"`
}

// Fill in conf struct from viper and extract all groups from `Remotes`
func (self *TypeConf) Init() {
	helper.Debug = Flag.Debug
	self.File = viper.ConfigFileUsed()
	viper.Unmarshal(&self)
	// Fill in ConfGroup
	for _, r := range self.Remotes {
		self.Groups.Add(&r.Group)
	}
	helper.ReportDebug(&Conf.File, "", true, true)
	helper.ReportDebug(&Flag, "Flag", true, false)
}

// Calculate remotes base on flag
func (self *TypeConf) MergeRemotes() {
	// Merge remote from flag "group"
	for _, g := range Flag.Groups {
		group := &g
		if self.Groups.Has(group) {
			self.MergedRemotes.AddArray(self.Remotes.GetByGroup(group))
		} else {
			log.Fatal("Group not in config: " + g)
			os.Exit(1)
		}
	}
	// Merge remote from flag "remote"
	for _, r := range Flag.Remotes {
		if self.Remotes.Has(&r) {
			self.MergedRemotes.Add(self.Remotes.GetByName(&r))
		} else {
			log.Fatal("Remote not in config: " + r)
			os.Exit(1)
		}
	}
	// If no remote is specified from command line, use all remotes in config
	if len(self.MergedRemotes) == 0 {
		self.MergedRemotes = append(self.MergedRemotes, self.Remotes...)
	}
	helper.ReportDebug(&Conf.MergedRemotes, "Merged Remote", true, false)
}
