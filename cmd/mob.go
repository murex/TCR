package cmd

import (
	"github.com/mengdaming/tcr/tcr"
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
		tcr.Start(tcr.Mob, toolchain, true)
	},
}

func init() {
	rootCmd.AddCommand(mobCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mobCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mobCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
