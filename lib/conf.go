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

	"github.com/J-Siu/go-helper/v2/basestruct"
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
	*basestruct.Base

	FileConf      string      `json:"FileConf"`
	Groups        Groups      `json:"Groups"`
	Remotes       Remotes     `json:"remotes"`
	Secrets       ConfSecrets `json:"secrets"`
	MergedRemotes Remotes     `json:"MergedRemotes"`
}

func (t *TypeConf) New(flagGroups, flagRemotes *[]string) {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeConf"
	prefix := t.MyType + ".New"

	t.setDefault()
	ezlog.Debug().N(prefix).N("Default").Lm(t).Out()

	t.readFileConf()
	if t.Err == nil {
		t.initGroups()
		t.mergeRemotes(flagGroups, flagRemotes)

		ezlog.Debug().N(prefix).N("Raw").Lm(t).Out()

		t.expand()

		ezlog.Debug().N(prefix).N("Expand").Lm(t).Out()
	}
}

func (t *TypeConf) readFileConf() {
	prefix := t.MyType + ".readFileConf"

	viper.SetConfigType("json")
	viper.SetConfigFile(file.TildeEnvExpand(t.FileConf))
	viper.AutomaticEnv() // read in environment variables that match
	t.Err = viper.ReadInConfig()

	if t.Err == nil {
		viper.Unmarshal(&t)
	} else {
		ezlog.Debug().N(prefix).M(t.Err).Out()
	}
}

// This should be called
//   - before reading config file
func (t *TypeConf) setDefault() {
	if t.FileConf == "" {
		t.FileConf = Default.FileConf
	}
}

// This should be called
//   - after reading config file
//   - before merging remote
func (t *TypeConf) initGroups() {
	for _, r := range t.Remotes {
		t.Groups.Add(&r.Group)
	}
}

// Calculate remotes base on flag
func (t *TypeConf) mergeRemotes(flagGroups, flagRemotes *[]string) {
	// Merge remote from flag "group"
	for _, g := range *flagGroups {
		group := &g
		if t.Groups.Has(group) {
			t.MergedRemotes.AddArray(t.Remotes.GetByGroup(group))
		} else {
			log.Fatal("Group not in config: " + g)
			os.Exit(1)
		}
	}
	// Merge remote from flag "remote"
	for _, r := range *flagRemotes {
		if t.Remotes.Has(&r) {
			t.MergedRemotes.Add(t.Remotes.GetByName(&r))
		} else {
			log.Fatal("Remote not in config: " + r)
			os.Exit(1)
		}
	}
	// If no remote is specified from command line, use all remotes in config
	if len(t.MergedRemotes) == 0 {
		t.MergedRemotes = append(t.MergedRemotes, t.Remotes...)
	}
	// helper.ReportDebug(&Conf.MergedRemotes, "Merged Remote", true, false)
}

func (t *TypeConf) expand() {
	t.FileConf = file.TildeEnvExpand(t.FileConf)
}
