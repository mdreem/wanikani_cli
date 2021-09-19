package cmd

import (
	"github.com/mdreem/wanikani_cli/levelprogress"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/cobra"
)

var levelProgressCmd = &cobra.Command{
	Use: "level_progress    Prints progress of the current level",
	Run: func(command *cobra.Command, args []string) {
		client := wanikani.CreateClient()
		levelprogress.PrintLevelProgress(client, command, args)
	},
}

func init() {
	RootCmd.AddCommand(levelProgressCmd)
}
