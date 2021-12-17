package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	Branch    = "main"
	Author    = "mobus"
	Email     = "<sv0202@163.com>"
	Date      = "2021-09-26"
	Commit    = "821288f"
	GoVersion = "go1.17.1 linux/amd64"
)

var version = &cli.Command{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "show version info",
	Action: func(c *cli.Context) error {
		fmt.Println("branch: ", Branch)
		fmt.Println("author: ", Author)
		fmt.Println("email: ", Email)
		fmt.Println("date: ", Date)
		fmt.Println("git commit: ", Commit)
		fmt.Println(GoVersion)
		return nil
	},
}
