package cmd

import (
	"github.com/mdreem/wanikani_cli/progressinfo"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/cobra"
)

var progressInfoCmd = &cobra.Command{
	Use: "progress_info     Computes overall progress",
	Run: func(command *cobra.Command, args []string) {
		client := wanikani.CreateClient()
		progressinfo.PrintProgressInfo(client, command, args)
	},
}

func init() {
	RootCmd.AddCommand(progressInfoCmd)
}
