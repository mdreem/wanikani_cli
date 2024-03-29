package main

import (
	"flag"
	"fmt"
	"github.com/mdreem/wanikani_cli/wanikani"
	"github.com/spf13/viper"
	"os"
)

func main() {
	commandInfo := parseArguments()

	initializeConfiguration()

	switch commandInfo.Command {
	case wanikani.LevelProgress:
		wanikani.PrintLevelProgress()
	case wanikani.UserInfo:
		wanikani.PrintUserInfo()
	case wanikani.ProgressInfo:
		wanikani.PrintProgressInfo()
	}
}

func parseArguments() wanikani.CommandInfo {
	if len(os.Args) < 2 {
		fmt.Println("expected a subcommand")
		fmt.Println("")
		fmt.Println("possible subcommands:")
		fmt.Println("")
		fmt.Println("level_progress    Prints progress of the current level")
		fmt.Println("user_info         Fetches user data")
		fmt.Println("progress_info     Computes overall progress")

		os.Exit(1)
	}

	levelProgress := flag.NewFlagSet("foo", flag.ExitOnError)
	userInfo := flag.NewFlagSet("foo", flag.ExitOnError)

	var command wanikani.Command

	switch os.Args[1] {
	case "level_progress":
		err := levelProgress.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Errorf("error parsing paramerers for level_progress: %v", err))
		}
		command = wanikani.LevelProgress
	case "user_info":
		err := userInfo.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Errorf("error parsing paramerers for user_info: %v", err))
		}
		command = wanikani.UserInfo
	case "progress_info":
		err := userInfo.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Errorf("error parsing paramerers for progress_info: %v", err))
		}
		command = wanikani.ProgressInfo
	default:
		fmt.Println("expected a subcommand")
		os.Exit(1)
	}

	return wanikani.CommandInfo{
		Command: command,
	}
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
