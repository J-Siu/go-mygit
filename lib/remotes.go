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

import "github.com/J-Siu/go-helper/v2/str"

// Array of Remote
//
// Embedded in TypeConf. No initialization
type Remotes []Remote

// Check array has Remote (by name)
func (remotes *Remotes) Has(name *string) bool {
	for _, r := range *remotes {
		if r.Name == *name {
			return true
		}
	}
	return false
}

// Add Remote
func (remotes *Remotes) Add(rP *Remote) {
	if rP != nil && !remotes.Has(&rP.Name) {
		*remotes = append(*remotes, *rP)
	}
}

// Add Remote Array
func (remotes *Remotes) AddArray(raP *Remotes) {
	if raP != nil {
		*remotes = append(*remotes, *raP...)
	}
}

// Get Remote by name
func (remotes *Remotes) GetByName(nameP *string) *Remote {
	for _, r := range *remotes {
		if r.Name == *nameP {
			return &r
		}
	}
	return nil
}

// Get all Remote in a group
func (remotes *Remotes) GetByGroup(groupP *string) *Remotes {
	var tmpRemotes Remotes
	for _, r := range *remotes {
		if r.Group == *groupP {
			tmpRemotes.Add(&r)
		}
	}
	return &tmpRemotes
}

// Get all Remote names
func (remotes *Remotes) GetNames() *[]string {
	var names []string
	for _, r := range *remotes {
		if !str.ArrayContains(&names, &r.Name) {
			names = append(names, r.Name)
		}
	}
	return &names
}
