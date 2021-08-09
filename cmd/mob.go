package cmd

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/tcr/engine"
	"github.com/mengdaming/tcr/tcr/runmode"
	"github.com/mengdaming/tcr/tcr/ui/cli"
	"github.com/spf13/cobra"
)

// mobCmd represents the mob command
var mobCmd = &cobra.Command{
	Use:   "mob",
	Short: "Run TCR in mob mode",
	Long: `
When used in "mob" mode, TCR ensures that any commit
is shared with other participants through calling git push-pull.`,
	Run: func(cmd *cobra.Command, args []string) {
		u := cli.New()
		params.Mode = runmode.Mob{}
		params.AutoPush = true
		params.PollingPeriod = tcr.DefaultPollingPeriod
		engine.Init(u, params)
	},
}

func init() {
	rootCmd.AddCommand(mobCmd)
}
