package cmd

import (
	"github.com/mengdaming/tcr/engine"
	"github.com/mengdaming/tcr/runmode"
	"github.com/mengdaming/tcr/ui/cli"
	"github.com/spf13/cobra"
)

// mobCmd represents the mob command
var mobCmd = &cobra.Command{
	Use:   "mob",
	Short: "Run TCR in mob mode",
	Long: `
When used in "mob" mode, TCR ensures that any commit
is shared with other participants through calling git push-pull.

This subcommand runs directly in the terminal (no GUI).`,
	Run: func(cmd *cobra.Command, args []string) {
		params.Mode = runmode.Mob{}
		params.AutoPush = params.Mode.AutoPushDefault()
		params.PollingPeriod = engine.DefaultPollingPeriod
		u := cli.New(params)
		u.Start(params.Mode)
	},
}

func init() {
	rootCmd.AddCommand(mobCmd)
}
