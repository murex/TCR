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

type paramValueBool struct {
	value        bool
	defaultValue bool
}

func (p *paramValueBool) reset() {
	p.value = p.defaultValue
}

// BoolParam is a parameter of type bool that can be handled by both viper and cobra frameworks
type BoolParam struct {
	s paramSettings
	v paramValueBool
}

func (param *BoolParam) addToCommand(cmd *cobra.Command) {
	flags := param.s.getCmdFlags(cmd)
	flags.BoolVarP(&param.v.value,
		param.s.cobraSettings.name,
		param.s.shorthand,
		false,
		param.s.usage)
	flag := flags.Lookup(param.s.cobraSettings.name)
	param.s.bindToViper(flag)
}

func (param *BoolParam) useDefaultValueIfNotSet() {
	const undefined bool = false
	// nolint:staticcheck // We want to use undefined variable explicitly here
	if param.v.value == undefined {
		if param.s.enabled {
			// nolint:staticcheck // We want to use undefined variable explicitly here
			if cfgValue := viper.GetBool(param.s.getViperKey()); cfgValue != undefined {
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
func (param *BoolParam) GetValue() bool {
	param.useDefaultValueIfNotSet()
	return param.v.value
}

func (param *BoolParam) reset() {
	param.v.reset()
	if param.s.enabled {
		viper.Set(param.s.getViperKey(), param.v.value)
	}
}
