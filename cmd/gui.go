package cmd

import (
	"github.com/mengdaming/tcr/engine"
	"github.com/mengdaming/tcr/runmode"
	"github.com/mengdaming/tcr/ui/gui"
	"github.com/spf13/cobra"
)

// mobCmd represents the mob command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch TCR GUI",
	Long: `
Run TCR application though a Graphical User Interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Mode = runmode.Mob{}
		params.AutoPush = params.Mode.AutoPushDefault()
		params.PollingPeriod = engine.DefaultPollingPeriod
		u := gui.New(params)
		u.Start(params.Mode)
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
