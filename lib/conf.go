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
	"log"
	"os"

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
func (c *TypeConf) Init(flag *TypeFlag) {
	c.File = viper.ConfigFileUsed()
	viper.Unmarshal(&c)
	// Fill in ConfGroup
	for _, r := range c.Remotes {
		c.Groups.Add(&r.Group)
	}
}

// Calculate remotes base on flag
func (c *TypeConf) MergeRemotes(flag *TypeFlag) {
	// Merge remote from flag "group"
	for _, g := range flag.Groups {
		group := &g
		if c.Groups.Has(group) {
			c.MergedRemotes.AddArray(c.Remotes.GetByGroup(group))
		} else {
			log.Fatal("Group not in config: " + g)
			os.Exit(1)
		}
	}
	// Merge remote from flag "remote"
	for _, r := range flag.Remotes {
		if c.Remotes.Has(&r) {
			c.MergedRemotes.Add(c.Remotes.GetByName(&r))
		} else {
			log.Fatal("Remote not in config: " + r)
			os.Exit(1)
		}
	}
	// If no remote is specified from command line, use all remotes in config
	if len(c.MergedRemotes) == 0 {
		c.MergedRemotes = append(c.MergedRemotes, c.Remotes...)
	}
}
