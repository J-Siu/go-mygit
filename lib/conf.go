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
- File, Group, MergeRemotes are filled in at runtime

- Remotes, Secrets are read from config file by viper
*/
type TypeConf struct {
	FileConf      string      `json:"FileConf"`
	Groups        Groups      `json:"Groups"`
	Remotes       Remotes     `json:"remotes"`
	Secrets       ConfSecrets `json:"secrets"`
	MergedRemotes Remotes     `json:"MergedRemotes"`
}

func (conf *TypeConf) Init() {
	prefix := "TypeConf.Init"

	conf.setDefault()

	helper.ReportDebug(conf.FileConf, prefix+": Config file", false, true)

	conf.readFileConf()

	conf.groups()
	conf.mergeRemotes()

	helper.ReportDebug(conf, prefix+": Raw", false, true)

	conf.expand()

	helper.ReportDebug(conf, prefix+": Expand", false, true)
}

func (conf *TypeConf) readFileConf() {
	prefix := "TypeConf.readFileConf"
	viper.SetConfigType("json")
	viper.SetConfigFile(helper.TildeEnvExpand(Conf.FileConf))
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&conf)
	} else {
		helper.Report(err.Error(), prefix, true, true)
		os.Exit(1)
	}
}

func (conf *TypeConf) setDefault() {
	if conf.FileConf == "" {
		conf.FileConf = Default.FileConf
	}
}

// This should be called
//   - after reading config file
//   - before merging remote
func (conf *TypeConf) groups() {
	for _, r := range conf.Remotes {
		conf.Groups.Add(&r.Group)
	}
}

// Calculate remotes base on flag
func (conf *TypeConf) mergeRemotes() {
	// Merge remote from flag "group"
	for _, g := range Flag.Groups {
		group := &g
		if conf.Groups.Has(group) {
			conf.MergedRemotes.AddArray(conf.Remotes.GetByGroup(group))
		} else {
			log.Fatal("Group not in config: " + g)
			os.Exit(1)
		}
	}
	// Merge remote from flag "remote"
	for _, r := range Flag.Remotes {
		if conf.Remotes.Has(&r) {
			conf.MergedRemotes.Add(conf.Remotes.GetByName(&r))
		} else {
			log.Fatal("Remote not in config: " + r)
			os.Exit(1)
		}
	}
	// If no remote is specified from command line, use all remotes in config
	if len(conf.MergedRemotes) == 0 {
		conf.MergedRemotes = append(conf.MergedRemotes, conf.Remotes...)
	}
	// helper.ReportDebug(&Conf.MergedRemotes, "Merged Remote", true, false)
}

func (conf *TypeConf) expand() {
	conf.FileConf = helper.TildeEnvExpand(conf.FileConf)
}
