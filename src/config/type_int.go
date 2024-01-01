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

package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type paramValueInt struct {
	value        int
	defaultValue int
}

func (p *paramValueInt) reset() {
	p.value = p.defaultValue
}

// IntParam is a parameter of type int that can be handled by both viper and cobra frameworks
type IntParam struct {
	s paramSettings
	v paramValueInt
}

func (param *IntParam) addToCommand(cmd *cobra.Command) {
	flags := param.s.getCmdFlags(cmd)
	flags.IntVarP(&param.v.value,
		param.s.cobraSettings.name,
		param.s.cobraSettings.shorthand,
		0,
		param.s.cobraSettings.usage)
	flag := flags.Lookup(param.s.cobraSettings.name)
	param.s.bindToViper(flag)
}

func (param *IntParam) useDefaultValueIfNotSet() {
	const undefined int = 0
	// nolint:gosimple // We want to use undefined variable explicitly here
	if param.v.value == undefined {
		if param.s.viperSettings.enabled {
			if cfgValue := viper.GetInt(param.s.getViperKey()); cfgValue != undefined {
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
func (param *IntParam) GetValue() int {
	param.useDefaultValueIfNotSet()
	return param.v.value
}

func (param *IntParam) reset() {
	param.v.reset()
	if param.s.enabled {
		viper.Set(param.s.getViperKey(), param.v.value)
	}
}
