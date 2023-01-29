package gojira

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type gojira struct {
	cli *cli.App
	ui  *UserInteface
}

func Run() {
	app.ui = &UserInteface{}
	app.cli = &cli.App{
		Name: "gojira",
		Usage: `quickly log time to jira/tempo through cli.

   Calling without arguments will try to detect issue from git branch, 
   otherwise it will display list of last updated issues you're are assigned to.`,
		Version: "0.3.1",
		Before: func(context *cli.Context) error {
			if context.Args().First() != "config" {
				// dont' check envs on ConfigCommand
				PrepareConfig()
			}
			return nil
		},
		Commands: []*cli.Command{
			LogWorkCommand,
			IssuesCommand,
			WorkLogsCommand,
			ConfigCommand,
			ViewIssueCommand,
		},
		Action: DefaultAction,
	}
	err := app.cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var app gojira
