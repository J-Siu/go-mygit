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

package set

import (
	"github.com/spf13/cobra"
)

var (
	tf flagsTF
)

// struct to setup True/False flags
type flagsTF struct {
	setFalse bool
	setTrue  bool
}

// initialize mutual exclusive and require one true/false flags for cobra command
func (t *flagsTF) initTrueFalse(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&t.setFalse, "false", "f", false, "false")
	cmd.Flags().BoolVarP(&t.setTrue, "true", "t", false, "true")
	cmd.MarkFlagsMutuallyExclusive("false", "true")
	cmd.MarkFlagsOneRequired("false", "true")
}

// initialize mutual exclusive and require one public/private flags for cobra command
// public: setTrue:=true
// public: setTrue:=false
func (t *flagsTF) initPublicPrivate(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&t.setFalse, "private", "", false, "private")
	cmd.Flags().BoolVarP(&t.setTrue, "public", "", false, "public")
	cmd.MarkFlagsMutuallyExclusive("private", "public")
	cmd.MarkFlagsOneRequired("private", "public")
}
