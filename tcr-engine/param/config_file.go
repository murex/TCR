package param

import (
	"github.com/spf13/cobra"
)

func NewConfigFileParam(cmd *cobra.Command) *StringParam {
	param := StringParam{
		s: paramSettings{
			viperSettings: viperSettings{
				enabled: false,
				keyPath: "",
				name:    "",
			},
			cobraSettings: cobraSettings{
				name:       "config",
				shorthand:  "c",
				usage:      "config file (default is $HOME/.tcr.yaml)",
				persistent: true,
			},
		},
		v: paramValueString{
			value:        "",
			defaultValue: "",
		},
	}
	param.addToCommand(cmd)
	return &param
}
