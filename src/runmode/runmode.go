/*
Copyright (c) 2022 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package runmode

// RunMode is the interface that any run mode needs to satisfy to bee used by the TCR engine
type RunMode interface {
	Name() string
	AutoPushDefault() bool
	IsMultiRole() bool
	IsInteractive() bool
	IsActive() bool
}

var (
	allModes = []RunMode{Mob{}, Solo{}, OneShot{}, Check{}, Log{}, Stats{}}
)

// InteractiveModes returns the list of names of available interactive run modes
func InteractiveModes() []string {
	var names []string
	for _, mode := range allModes {
		if mode.IsInteractive() {
			names = append(names, mode.Name())
		}
	}
	return names
}

// Map returns the list of available run modes as a map of strings
func Map() map[string]RunMode {
	var m = make(map[string]RunMode)
	for _, mode := range allModes {
		m[mode.Name()] = mode
	}
	return m
}
