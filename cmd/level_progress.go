package cmd

import (
	"github.com/mdreem/wanikani_cli/levelprogress"
	"github.com/spf13/cobra"
)

var levelProgressCmd = &cobra.Command{
	Use: "level_progress    Prints progress of the current level",
	Run: func(c *cobra.Command, args []string) {
		levelprogress.PrintLevelProgress(c, args)
	},
}

func init() {
	RootCmd.AddCommand(levelProgressCmd)
}
