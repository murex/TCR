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

func (s viperSettings) getViperKey() string {
	if s.enabled {
		return fmt.Sprintf("%v.%v", s.keyPath, s.name)
	} else {
		return ""
	}
}

func (s viperSettings) bindToViper(flag *pflag.Flag) {
	if s.enabled {
		_ = viper.BindPFlag(s.getViperKey(), flag)
	}
}

func (c cobraSettings) getCmdFlags(cmd *cobra.Command) *pflag.FlagSet {
	if c.persistent {
		return cmd.PersistentFlags()
	} else {
		return cmd.Flags()
	}
}
