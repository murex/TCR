package param

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type paramValueString struct {
	value        string
	defaultValue string
}

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
				param.v.value = param.v.defaultValue
			}
		} else {
			param.v.value = param.v.defaultValue
		}
	}
}

func (param *StringParam) GetValue() string {
	param.useDefaultValueIfNotSet()
	return param.v.value
}
