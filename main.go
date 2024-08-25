package main

import (
	"fmt"
	"github-stars/commands"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	_ = godotenv.Load()

	app := cli.App{
		Name:        "gsv",
		Usage:       "查找并展示 github 用户 stars repos",
		Description: "查找并展示 github 用户 stars repos",
		Commands: cli.Commands{
			commands.AllRepos(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("cli exec command error: %s", err.Error())
	}
}
