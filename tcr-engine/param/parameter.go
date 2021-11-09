package param

import (
	"fmt"
	"github.com/spf13/cobra"
)

type viperParameter interface {
	getViperKey() string
}

type cobraParameter interface {
	Init(cmd *cobra.Command)
}

type Parameter interface {
	cobraParameter
	viperParameter
	DefaultWhenUndefined()
}

type viperSettings struct {
	category string
	name     string
}

type cobraSettings struct {
	Name      string
	shorthand string
	usage     string
}

type Settings struct {
	viperSettings
	cobraSettings
}

func (s viperSettings) getViperKey() string {
	return fmt.Sprintf("%v.%v", s.category, s.name)
}
