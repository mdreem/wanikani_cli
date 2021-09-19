package cmd

import (
	"github.com/mdreem/wanikani_cli/level_progress"
	"github.com/spf13/cobra"
)

var levelProgressCmd = &cobra.Command{
	Use: "level_progress    Prints progress of the current level",
	Run: func(c *cobra.Command, args []string) {
		level_progress.PrintLevelProgress(c, args)
	},
}

func init() {
	RootCmd.AddCommand(levelProgressCmd)
}
