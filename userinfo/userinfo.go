package userinfo

import (
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/cobra"
)

func PrintUserInfo(client wanikani.WanikaniClient, _ *cobra.Command, _ []string) {
	userInformation := client.FetchUserInformation()

	fmt.Printf("Username: %s\n", userInformation.Username)
	fmt.Printf("Level: %s\n", userInformation.Level)
	fmt.Printf("Profile URL: %s\n", userInformation.ProfileURL)

	fmt.Printf("CurrentVacationStartedAt: %s\n", userInformation.CurrentVacationStartedAt)
	fmt.Println("Subscription:")
	fmt.Printf("\tActive: %t\n", userInformation.Subscription.Active)
	fmt.Printf("\tMax level granted: %s\n", userInformation.Subscription.MaxLevelGranted.String())
	fmt.Printf("\tPeriod ends at: %s\n", userInformation.Subscription.PeriodEndsAt)
	fmt.Printf("\tSubscription type: %s\n", userInformation.Subscription.Type)
}
