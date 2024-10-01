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

// AddVariantParam adds variant parameter to the provided command
func AddVariantParam(cmd *cobra.Command) *StringParam {
	param := StringParam{
		s: paramSettings{
			viperSettings: viperSettings{
				enabled: true,
				keyPath: "config.tcr",
				name:    "variant",
			},
			cobraSettings: cobraSettings{
				name:       "variant",
				shorthand:  "r",
				usage:      "indicate the variant to be used by TCR: relaxed (default) or btcr",
				persistent: true,
			},
		},
		v: paramValueString{
			value:        "",
			defaultValue: "relaxed",
		},
	}
	param.addToCommand(cmd)
	return &param
}
