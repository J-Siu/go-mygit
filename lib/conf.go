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

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/spf13/viper"
)

var Default = TypeConf{
	FileConf: "$HOME/.config/go-mygit.json",
}

/*
- File, Group, MergeRemotes are filled in at runtime

- Remotes, Secrets are read from config file by viper
*/
type TypeConf struct {
	Err    error
	myType string
	init   bool

	FileConf      string      `json:"FileConf"`
	Groups        Groups      `json:"Groups"`
	Remotes       Remotes     `json:"remotes"`
	Secrets       ConfSecrets `json:"secrets"`
	MergedRemotes Remotes     `json:"MergedRemotes"`
}

func (c *TypeConf) New(flagGroups, flagRemotes *[]string) {
	c.init = true
	c.myType = "TypeConf"
	prefix := c.myType + ".Init"

	c.setDefault()
	ezlog.Debug().N(prefix).Nn("Default").M(c).Out()

	c.readFileConf()
	if c.Err == nil {
		c.initGroups()
		c.mergeRemotes(flagGroups, flagRemotes)

		ezlog.Debug().N(prefix).Nn("Raw").M(c).Out()

		c.expand()

		ezlog.Debug().N(prefix).Nn("Expand").M(c).Out()
	}
}

func (c *TypeConf) readFileConf() {
	prefix := c.myType + ".readFileConf"

	viper.SetConfigType("json")
	viper.SetConfigFile(file.TildeEnvExpand(c.FileConf))
	viper.AutomaticEnv() // read in environment variables that match
	c.Err = viper.ReadInConfig()

	if c.Err == nil {
		viper.Unmarshal(&c)
	} else {
		ezlog.Debug().N(prefix).M(c.Err).Out()
	}
}

// This should be called
//   - before reading config file
func (c *TypeConf) setDefault() {
	if c.FileConf == "" {
		c.FileConf = Default.FileConf
	}
}

// This should be called
//   - after reading config file
//   - before merging remote
func (c *TypeConf) initGroups() {
	for _, r := range c.Remotes {
		c.Groups.Add(&r.Group)
	}
}

// Calculate remotes base on flag
func (c *TypeConf) mergeRemotes(flagGroups, flagRemotes *[]string) {
	// Merge remote from flag "group"
	for _, g := range *flagGroups {
		group := &g
		if c.Groups.Has(group) {
			c.MergedRemotes.AddArray(c.Remotes.GetByGroup(group))
		} else {
			log.Fatal("Group not in config: " + g)
			os.Exit(1)
		}
	}
	// Merge remote from flag "remote"
	for _, r := range *flagRemotes {
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
	// helper.ReportDebug(&Conf.MergedRemotes, "Merged Remote", true, false)
}

func (c *TypeConf) expand() {
	c.FileConf = file.TildeEnvExpand(c.FileConf)
}
