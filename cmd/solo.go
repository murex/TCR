package cmd

import (
	"github.com/mengdaming/tcr/tcr"
	"github.com/mengdaming/tcr/trace"
	"github.com/spf13/cobra"
)

// soloCmd represents the solo command
var soloCmd = &cobra.Command{
	Use:   "solo",
	Short: "Run TCR in solo mode",
	Long: `
When used in "solo" mode, TCR only commits changes locally.
It never pushes or pulls to a remote repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		trace.HorizontalLine()
		trace.Info("Running in solo mode")
		tcr.Start()
	},
}

func init() {
	rootCmd.AddCommand(soloCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// soloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// soloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
