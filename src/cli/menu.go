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

package cli

import (
	"bytes"
	"errors"
	"fmt"
)

const menuArrow = "─▶"

type (
	menuAction  func()
	menuEnabler func() bool

	menuOption struct {
		shortcut    byte
		description string
		enabler     menuEnabler
		action      menuAction
		quitOption  bool
	}

	menu struct {
		title   string
		options []*menuOption
	}
)

func newMenuOption(s byte, d string, e menuEnabler, a menuAction, q bool) *menuOption { //nolint:revive
	return &menuOption{
		shortcut:    s,
		description: d,
		enabler:     e,
		action:      a,
		quitOption:  q,
	}
}

func (mo *menuOption) getDescription() string {
	return mo.description
}

func (mo *menuOption) getShortcut() byte {
	return mo.shortcut
}

func (mo *menuOption) run() error {
	if mo.action == nil {
		return errors.New("menu option action is not set")
	}
	mo.action()
	return nil
}

func (mo *menuOption) matchShortcut(b byte) bool {
	return bytes.ToUpper([]byte{b})[0] == mo.getShortcut()
}

func (mo *menuOption) isQuitOption() bool {
	return mo.quitOption
}

func (mo *menuOption) isEnabled() bool {
	return mo.enabler == nil || mo.enabler()
}

func (mo *menuOption) toString() string {
	return fmt.Sprintf("\t%c %s %s", mo.getShortcut(), menuArrow, mo.getDescription())
}

func newMenu(title string) *menu {
	return &menu{
		title: title,
	}
}

func (m *menu) getTitle() string {
	return m.title
}

func (m *menu) setTitle(title string) {
	m.title = title
}

func (m *menu) addOptions(options ...*menuOption) {
	m.options = append(m.options, options...)
}

func (m *menu) getOptions() (out []*menuOption) {
	for _, option := range m.options {
		if option.isEnabled() {
			out = append(out, option)
		}
	}
	return out
}

func (m *menu) matchAndRun(input byte) (matched bool, quit bool) {
	for _, option := range m.getOptions() {
		if !matched && option.matchShortcut(input) {
			matched = true
			_ = option.run()
			if option.quitOption {
				return matched, true
			}
		}
	}
	return matched, false
}
