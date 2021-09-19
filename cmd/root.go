package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var RootCmd = &cobra.Command{
	Use: "wanikani_cli",
	Run: runCommand,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("could not execute command: %v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		initializeConfiguration()
	})
}

func runCommand(_ *cobra.Command, _ []string) {
}

func initializeConfiguration() {
	viper.SetConfigName("wanikani")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error: %s", err))
	}
}
