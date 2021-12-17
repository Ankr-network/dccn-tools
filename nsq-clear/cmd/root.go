package cmd

import (
	"github.com/Ankr-network/dccn-tools/nsq-clear/app"
	"github.com/Ankr-network/dccn-tools/nsq-clear/share"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func MainRun() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     share.NsqAdmin,
			Value:    "http://127.0.0.1:4171",
			Required: true,
			Usage:    "nsq admin endpoint",
		},
		&cli.Uint64Flag{
			Name:  share.Threshold,
			Value: 100000,
			Usage: "nsq queue threshold， empty when exceeds",
		},
		&cli.DurationFlag{
			Name:  share.Schedule,
			Value: time.Minute * 5,
			Usage: "how often does it check",
		},
		&cli.StringFlag{
			Name:  share.LogLevel,
			Value: "info",
			Usage: "log level",
		},
	}

	svr := cli.NewApp()
	svr.Action = app.MainServe
	svr.Flags = flags

	appendCmdList(svr, version)
	err := svr.Run(os.Args)
	if err != nil {
		log.Fatal("Service Crash ", err)
	}
}

func appendCmdList(app *cli.App, subcmd *cli.Command) {
	app.Commands = append(app.Commands, subcmd)
}
