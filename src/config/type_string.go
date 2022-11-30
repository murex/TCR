/*
Copyright (c) 2021 Murex

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
	"github.com/spf13/viper"
)

type paramValueString struct {
	value        string
	defaultValue string
}

func (p *paramValueString) reset() {
	p.value = p.defaultValue
}

// StringParam is a parameter of type string that can be handled by both viper and cobra frameworks
type StringParam struct {
	s paramSettings
	v paramValueString
}

func (param *StringParam) addToCommand(cmd *cobra.Command) {
	flags := param.s.getCmdFlags(cmd)
	flags.StringVarP(&param.v.value,
		param.s.cobraSettings.name,
		param.s.cobraSettings.shorthand,
		"",
		param.s.cobraSettings.usage)

	flag := flags.Lookup(param.s.cobraSettings.name)
	param.s.bindToViper(flag)
}

func (param *StringParam) useDefaultValueIfNotSet() {
	const undefined string = ""
	if param.v.value == undefined {
		if param.s.viperSettings.enabled {
			if cfgValue := viper.GetString(param.s.getViperKey()); cfgValue != undefined {
				param.v.value = cfgValue
			} else {
				param.reset()
			}
		} else {
			param.reset()
		}
	}
}

// GetValue returns the current value for this parameter
func (param *StringParam) GetValue() string {
	param.useDefaultValueIfNotSet()
	return param.v.value
}

func (param *StringParam) reset() {
	param.v.reset()
	if param.s.enabled {
		viper.Set(param.s.getViperKey(), param.v.value)
	}
}
