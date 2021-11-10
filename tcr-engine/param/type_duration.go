package param

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

type paramValueDuration struct {
	value        time.Duration
	defaultValue time.Duration
}

type DurationParam struct {
	s paramSettings
	v paramValueDuration
}

func (param *DurationParam) addToCommand(cmd *cobra.Command) {
	flags := param.s.getCmdFlags(cmd)
	flags.DurationVarP(&param.v.value,
		param.s.cobraSettings.name,
		param.s.cobraSettings.shorthand,
		0,
		param.s.cobraSettings.usage)
	flag := flags.Lookup(param.s.cobraSettings.name)
	param.s.bindToViper(flag)
}

func (param *DurationParam) useDefaultValueIfNotSet() {
	const undefined time.Duration = 0
	if param.v.value == undefined {
		if param.s.viperSettings.enabled {
			if cfgValue := viper.GetDuration(param.s.getViperKey()); cfgValue != undefined {
				param.v.value = cfgValue
			} else {
				param.v.value = param.v.defaultValue
			}
		} else {
			param.v.value = param.v.defaultValue
		}
	}
}

func (param *DurationParam) GetValue() time.Duration {
	param.useDefaultValueIfNotSet()
	return param.v.value
}
