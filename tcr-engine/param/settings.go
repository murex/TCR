package param

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
	} else {
		return ""
	}
}

func (vs viperSettings) bindToViper(flag *pflag.Flag) {
	if vs.enabled {
		_ = viper.BindPFlag(vs.getViperKey(), flag)
	}
}

func (cs cobraSettings) getCmdFlags(cmd *cobra.Command) *pflag.FlagSet {
	if cs.persistent {
		return cmd.PersistentFlags()
	} else {
		return cmd.Flags()
	}
}
