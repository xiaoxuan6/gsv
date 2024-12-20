package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/xiaoxuan6/gsv/commands"
	"os"
)

func main() {
	_ = godotenv.Load()

	app := cli.App{
		Name:        "gsv",
		Usage:       "查找或下载 github 用户 stars repos",
		Description: "查找或下载 github 用户 stars repos",
		Commands: cli.Commands{
			commands.Search(),
			commands.Download(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("cli exec command error: %s", err.Error())
	}
}
