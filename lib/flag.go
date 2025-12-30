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

// Holding all flags from command line
type TypeFlag struct {
	Debug           bool     // Enable debug output
	Groups          []string // Groups specified in command line
	NoParallel      bool     // Do not process in parallel(go routine)
	NoSkip          bool     // Flag for not skipping empty output
	NoTitle         bool     // Do not print title in output
	Page            int      // Page number of repository listing
	PushAll         bool     // Flag for git push
	RemoteRemoveAll bool     // Flag for git remote remove all
	Remotes         []string // Remotes specified in command line
	Secret          Secret   // Secret specified in command line
	SecretsDel      []string // Secrets specified in command line
	Tag             bool     // Flag for git push, pull

	StatusOnly bool
	SingleLine bool
}
