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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type viperParameter interface {
	getViperKey() string
	bindToViper(flag *pflag.Flag)
}

type cobraParameter interface {
	attachToCommand(cmd *cobra.Command)
	getCmdFlags(cmd *cobra.Command) *pflag.FlagSet
}

// Parameter defines a parameter that can be handled by both viper and cobra frameworks
type Parameter interface {
	cobraParameter
	viperParameter
	useDefaultValueIfNotSet()
}

type viperSettings struct {
	enabled bool
	keyPath string
	name    string
}

type cobraSettings struct {
	name       string
	shorthand  string
	usage      string
	persistent bool
}

type paramSettings struct {
	viperSettings
	cobraSettings
}

func (vs viperSettings) getViperKey() string {
	if vs.enabled {
		return fmt.Sprintf("%v.%v", vs.keyPath, vs.name)
	}
	return ""
}

func (vs viperSettings) bindToViper(flag *pflag.Flag) {
	if vs.enabled {
		_ = viper.BindPFlag(vs.getViperKey(), flag)
	}
}

func (cs cobraSettings) getCmdFlags(cmd *cobra.Command) *pflag.FlagSet {
	if cs.persistent {
		return cmd.PersistentFlags()
	}
	return cmd.Flags()
}
