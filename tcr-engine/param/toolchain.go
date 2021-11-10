package param

import (
	"github.com/spf13/cobra"
)

func NewToolchainParam(cmd *cobra.Command) *StringParam {
	param := StringParam{
		s: paramSettings{
			viperSettings: viperSettings{
				enabled: true,
				keyPath: "params.tcr",
				name:    "toolchain",
			},
			cobraSettings: cobraSettings{
				name:       "toolchain",
				shorthand:  "t",
				usage:      "indicate the toolchain to be used by TCR",
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
