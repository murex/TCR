package param

import (
	"github.com/murex/tcr/tcr-engine/settings"
	"github.com/spf13/cobra"
)

func NewPollingPeriodParam(cmd *cobra.Command) *DurationParam {
	param := DurationParam{
		s: paramSettings{
			viperSettings: viperSettings{
				enabled: true,
				keyPath: "params.git",
				name:    "polling-period",
			},
			cobraSettings: cobraSettings{
				name:       "polling",
				shorthand:  "o",
				usage:      "set git polling period when running as navigator",
				persistent: true,
			},
		},
		v: paramValueDuration{
			value:        0,
			defaultValue: settings.DefaultPollingPeriod,
		},
	}
	param.addToCommand(cmd)
	return &param
}
