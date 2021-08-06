package cmd

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/runmode"
	"github.com/mengdaming/tcr/tcr/ui/gui"

	"github.com/spf13/cobra"
)

// mobCmd represents the mob command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch TCR GUI",
	Long: `
Runs TCR application though a Graphical User Interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		u := gui.New()
		params.Mode = runmode.Mob{}
		params.AutoPush = true
		params.PollingPeriod = tcr.DefaultPollingPeriod
		engine.Init(u, params)
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
