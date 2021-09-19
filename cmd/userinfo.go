package cmd

import (
	"github.com/mdreem/wanikani_cli/userinfo"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/cobra"
)

var userinfoCmd = &cobra.Command{
	Use: "user_info         Fetches user data",
	Run: func(command *cobra.Command, args []string) {
		client := wanikani.CreateClient()
		userinfo.PrintUserInfo(client, command, args)
	},
}

func init() {
	RootCmd.AddCommand(userinfoCmd)
}
