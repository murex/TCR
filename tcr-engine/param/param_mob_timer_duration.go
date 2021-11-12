package param

import (
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/spf13/cobra"
)

func NewMobTimerDurationParam(cmd *cobra.Command) *DurationParam {
	param := DurationParam{
		s: paramSettings{
			viperSettings: viperSettings{
				enabled: true,
				keyPath: "params.timer",
				name:    "duration",
			},
			cobraSettings: cobraSettings{
				name:       "duration",
				shorthand:  "d",
				usage:      "set the duration for role rotation countdown timer",
				persistent: true,
			},
		},
		v: paramValueDuration{
			value:        0,
			defaultValue: settings.DefaultMobTurnDuration,
		},
	}
	param.addToCommand(cmd)
	return &param
}
