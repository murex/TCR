package param

import (
	"github.com/spf13/cobra"
)

func NewAutoPushParam(cmd *cobra.Command) *BoolParam {
	param := BoolParam{
		s: paramSettings{
			viperSettings: viperSettings{
				enabled: true,
				keyPath: "params.git",
				name:    "auto-push",
			},
			cobraSettings: cobraSettings{
				name:       "auto-push",
				shorthand:  "p",
				usage:      "enable git push after every commit",
				persistent: true,
			},
		},
		v: paramValueBool{
			value:        false,
			defaultValue: false,
		},
	}
	param.addToCommand(cmd)
	return &param
}
