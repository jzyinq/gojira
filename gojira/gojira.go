package gojira

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

type gojira struct {
	cli            *cli.App
	ui             *UserInteface
	time           *time.Time
	workLogs       Worklogs
	workLogsIssues WorklogsIssues
}

func Run() {
	// Open the log file
	logFile, err := os.OpenFile("/tmp/gojira.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		logrus.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	// Set the log output to write to the file
	logrus.SetOutput(logFile)
	// Now log messages will be written to the file
	appTimer := time.Now().UTC()
	logrus.Infof("gojira version %s started", projectVersion)
	logrus.Infof("current time %s", appTimer)
	app.ui = &UserInteface{}
	app.time = &appTimer
	app.cli = &cli.App{
		Name: "gojira",
		Usage: `quickly log time to jira/tempo through cli.

   Calling without arguments will try to detect issue from git branch, 
   otherwise it will display list of last updated issues you're are assigned to.`,
		Version: projectVersion,
		Before: func(context *cli.Context) error {
			if context.Args().First() != "config" {
				// dont' check envs on ConfigCommand
				PrepareConfig()
			}
			if context.IsSet("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "Enable debug log level",
			},
		},
		Commands: []*cli.Command{
			LogWorkCommand,
			IssuesCommand,
			WorklogsCommand,
			ConfigCommand,
			ViewIssueCommand,
		},
		Action: DefaultAction,
	}
	err = app.cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var app gojira
