package param

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

type MobTimerDurationParam struct {
	settings     Settings
	defaultValue time.Duration
	value        time.Duration
}

func New() *MobTimerDurationParam {
	return &MobTimerDurationParam{
		settings: Settings{
			viperSettings: viperSettings{
				category: "params",
				name:     "duration",
			},
			cobraSettings: cobraSettings{
				Name:      "duration",
				shorthand: "d",
				usage:     "set the duration for role rotation countdown timer",
			},
		},
		defaultValue: 0,
		value:        0,
	}
}

func (param *MobTimerDurationParam) Init(cmd *cobra.Command) {

	cmd.PersistentFlags().DurationVarP(&param.value,
		param.settings.cobraSettings.Name,
		param.settings.cobraSettings.shorthand,
		0,
		param.settings.cobraSettings.usage)
	_ = viper.BindPFlag(param.settings.getViperKey(), cmd.PersistentFlags().Lookup(param.settings.cobraSettings.Name))
}

func (param *MobTimerDurationParam) DefaultWhenUndefined() {
	panic("implement me")
}
