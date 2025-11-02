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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_menu_option_get_description(t *testing.T) {
	mo := newMenuOption(0, "menu option description", nil, nil, false)
	assert.Equal(t, "menu option description", mo.getDescription())
}

func Test_menu_option_get_shortcut(t *testing.T) {
	mo := newMenuOption('D', "", nil, nil, false)
	assert.Equal(t, byte('D'), mo.getShortcut())
}

func Test_menu_option_to_string(t *testing.T) {
	mo := newMenuOption('Z', "menu option description", nil, nil, false)
	assert.Equal(t, "\tZ "+menuArrow+" menu option description", mo.toString())
}

func Test_menu_option_action(t *testing.T) {
	var result string
	mo := newMenuOption(0, "", nil, func() {
		result = "some value"
	}, false)
	err := mo.run()
	assert.NoError(t, err)
	assert.Equal(t, "some value", result)
}

func Test_menu_option_run_returns_when_no_action_is_set(t *testing.T) {
	mo := newMenuOption(0, "", nil, nil, false)
	assert.Error(t, mo.run())
}

func Test_menu_get_title(t *testing.T) {
	m := newMenu("menu title")
	assert.Equal(t, "menu title", m.getTitle())
}

func Test_menu_set_title(t *testing.T) {
	m := newMenu("some title")
	m.setTitle("new title")
	assert.Equal(t, "new title", m.getTitle())
}

func Test_menu_add_one_option(t *testing.T) {
	m := newMenu("menu title")
	mo := newMenuOption('X', "some description", nil, nil, false)
	assert.Equal(t, 0, len(m.options))
	m.addOptions(mo)
	assert.Equal(t, 1, len(m.options))
	assert.Equal(t, m.options[0], mo)
}

func Test_menu_add_multiple_options(t *testing.T) {
	m := newMenu("menu title")
	mo1 := newMenuOption('X', "some description", nil, nil, false)
	mo2 := newMenuOption('Y', "some description", nil, nil, false)
	assert.Equal(t, 0, len(m.options))
	m.addOptions(mo1, mo2)
	assert.Equal(t, 2, len(m.options))
	assert.Equal(t, m.options[0], mo1)
	assert.Equal(t, m.options[1], mo2)
}

func Test_menu_get_options(t *testing.T) {
	m := newMenu("menu title")
	mo1 := newMenuOption('X', "some description", nil, nil, false)
	m.addOptions(mo1)
	mo2 := newMenuOption('Y', "some description", nil, nil, false)
	m.addOptions(mo2)
	assert.Equal(t, m.getOptions(), []*menuOption{mo1, mo2})
}

func Test_menu_option_quit_option(t *testing.T) {
	tests := []struct {
		desc     string
		option   bool
		expected bool
	}{
		{"quit option set", true, true},
		{"quit option not set", false, false},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			mo := newMenuOption(0, "", nil, nil, test.option)
			assert.Equal(t, test.expected, mo.isQuitOption())
		})
	}
}

func Test_menu_option_match_shortcut(t *testing.T) {
	tests := []struct {
		desc     string
		shortcut byte
		input    byte
		expected bool
	}{
		{"X shortcut with X input", 'X', 'X', true},
		{"X shortcut with x input", 'X', 'x', true},
		{"Y shortcut with X input", 'Y', 'X', false},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			mo := newMenuOption(test.shortcut, "", nil, nil, false)
			assert.Equal(t, test.expected, mo.matchShortcut(test.input))
		})
	}
}

func Test_enable_menu_option(t *testing.T) {
	tests := []struct {
		desc     string
		enable   bool
		expected bool
	}{
		{"enable menu option", true, true},
		{"disable menu option", false, false},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			mo := newMenuOption(0, "", func() bool {
				return test.enable
			}, nil, false)
			assert.Equal(t, test.expected, mo.isEnabled())
		})
	}
}

func Test_menu_options_are_enabled_by_default(t *testing.T) {
	mo := newMenuOption(0, "", nil, nil, false)
	assert.True(t, mo.isEnabled())
}

func Test_menu_get_options_drops_disabled_options(t *testing.T) {
	m := newMenu("menu title")
	mo1 := newMenuOption('X', "some description", func() bool {
		return false
	}, nil, false)
	m.addOptions(mo1)
	mo2 := newMenuOption('Y', "some description", func() bool {
		return true
	}, nil, false)
	m.addOptions(mo2)
	assert.Equal(t, m.getOptions(), []*menuOption{mo2})
}

func Test_menu_match_and_run(t *testing.T) {
	tests := []struct {
		desc          string
		enabledFlag   bool
		quitFlag      bool
		input         byte
		expectedMatch bool
		expectedRun   bool
		expectedQuit  bool
	}{
		{
			"enabled regular option with matching input", true, false, 'X',
			true, true, false,
		},
		{
			"enabled regular option with non-matching input", true, false, 'Y',
			false, false, false,
		},
		{
			"disabled regular option with matching input", false, false, 'X',
			false, false, false,
		},
		{
			"disabled regular option with non-matching input", false, false, 'Y',
			false, false, false,
		},
		{
			"enabled quit option with matching input", true, true, 'X',
			true, true, true,
		},
		{
			"enabled quit option with non-matching input", true, true, 'Y',
			false, false, false,
		},
		{
			"disabled quit option with matching input", false, true, 'X',
			false, false, false,
		},
		{
			"disabled quit option with non-matching input", false, true, 'Y',
			false, false, false,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			m := newMenu("")
			var didRun bool
			m.addOptions(newMenuOption('X', "",
				func() bool {
					return test.enabledFlag
				},
				func() {
					didRun = true
				}, test.quitFlag))
			didMatch, shouldQuit := m.matchAndRun(test.input)
			assert.Equal(t, test.expectedMatch, didMatch, "input matching")
			assert.Equal(t, test.expectedRun, didRun, "action triggering")
			assert.Equal(t, test.expectedQuit, shouldQuit, "menu quit flag")
		})
	}
}
