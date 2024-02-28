/*
Copyright (c) 2024 Murex

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

package role

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Role provides the interface that a role must implement in order to be used by TCR engine
type Role interface {
	Name() string
	LongName() string
	RunsWithTimer() bool
}

func longName(r Role) string {
	return cases.Title(language.English).String(r.Name()) + " role"
}

var allRoles = []Role{Driver{}, Navigator{}}

// All returns the list of existing roles in TCR
func All() []Role {
	return allRoles
}

// FromName returns a role instance based on its name
func FromName(name string) Role {
	for _, r := range allRoles {
		if name == r.Name() {
			return r
		}
	}
	return nil
}
