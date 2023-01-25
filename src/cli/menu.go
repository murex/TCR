/*
Copyright (c) 2023 Murex

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
)

type menuAction func()

type menuOption struct {
	shortcut    byte
	description string
	help        string
	action      menuAction
	quitOption  bool
	enabled     bool
}

func newMenuOption(s byte, d string, h string, a menuAction, q bool) *menuOption { //nolint:revive
	return &menuOption{
		shortcut:    s,
		description: d,
		help:        h,
		action:      a,
		quitOption:  q,
		enabled:     true,
	}
}

func (mo *menuOption) getDescription() string {
	return mo.description
}

func (mo *menuOption) getShortcut() byte {
	return mo.shortcut
}

func (mo *menuOption) getHelp() string {
	return mo.help
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

func (mo *menuOption) setEnabled(enable bool) {
	mo.enabled = enable
}

func (mo *menuOption) isEnabled() bool {
	return mo.enabled
}

type menu struct {
	title   string
	options []*menuOption
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

func (m *menu) getOptions() []*menuOption {
	return m.options
}
