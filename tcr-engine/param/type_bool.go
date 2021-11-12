package param

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type paramValueBool struct {
	value        bool
	defaultValue bool
}

type BoolParam struct {
	s paramSettings
	v paramValueBool
}

func (param *BoolParam) addToCommand(cmd *cobra.Command) {
	flags := param.s.getCmdFlags(cmd)
	flags.BoolVarP(&param.v.value,
		param.s.cobraSettings.name,
		param.s.cobraSettings.shorthand,
		false,
		param.s.cobraSettings.usage)
	flag := flags.Lookup(param.s.cobraSettings.name)
	param.s.bindToViper(flag)
}

func (param *BoolParam) useDefaultValueIfNotSet() {
	const undefined bool = false
	if param.v.value == undefined {
		if param.s.viperSettings.enabled {
			if cfgValue := viper.GetBool(param.s.getViperKey()); cfgValue != undefined {
				param.v.value = cfgValue
			} else {
				param.v.value = param.v.defaultValue
			}
		} else {
			param.v.value = param.v.defaultValue
		}
	}
}

func (param *BoolParam) GetValue() bool {
	param.useDefaultValueIfNotSet()
	return param.v.value
}
