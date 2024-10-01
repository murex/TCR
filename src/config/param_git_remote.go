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

package config

import (
	"github.com/spf13/cobra"
)

// AddGitRemoteParam adds a parameter allowing to specify the name of the git remote repository
func AddGitRemoteParam(cmd *cobra.Command) *StringParam {
	param := StringParam{
		s: paramSettings{
			viperSettings: viperSettings{
				enabled: false,
				keyPath: "",
				name:    "",
			},
			cobraSettings: cobraSettings{
				name:       "git-remote",
				shorthand:  "g",
				usage:      "name of the git remote repository to sync with (default: \"origin\")",
				persistent: true,
			},
		},
		v: paramValueString{
			value:        "",
			defaultValue: "origin",
		},
	}
	param.addToCommand(cmd)
	return &param
}
