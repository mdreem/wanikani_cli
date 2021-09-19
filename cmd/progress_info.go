package cmd

import (
	"github.com/mdreem/wanikani_cli/progress_info"
	"github.com/spf13/cobra"
)

var progressInfoCmd = &cobra.Command{
	Use: "progress_info     Computes overall progress",
	Run: func(c *cobra.Command, args []string) {
		progress_info.PrintProgressInfo(c, args)
	},
}

func init() {
	RootCmd.AddCommand(progressInfoCmd)
}
