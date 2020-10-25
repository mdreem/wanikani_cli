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
		fmt.Println("Not implemented yet")
	}
}

func parseArguments() wanikani.CommandInfo {
	if len(os.Args) < 2 {
		fmt.Println("expected a subcommand")
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
		fmt.Println("subcommand 'level_progress'")
		command = wanikani.LevelProgress
	case "user_info":
		err := userInfo.Parse(os.Args[2:])
		if err != nil {
			panic(fmt.Errorf("error parsing paramerers for user_info: %v", err))
		}
		fmt.Println("subcommand 'userInfo'")
		command = wanikani.UserInfo
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
