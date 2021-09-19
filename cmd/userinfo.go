package cmd

import (
	"github.com/mdreem/wanikani_cli/userinfo"
	"github.com/spf13/cobra"
)

var userinfoCmd = &cobra.Command{
	Use: "user_info         Fetches user data",
	Run: func(c *cobra.Command, args []string) {
		userinfo.PrintUserInfo(c, args)
	},
}

func init() {
	RootCmd.AddCommand(userinfoCmd)
}
