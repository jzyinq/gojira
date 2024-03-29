package gojira

import (
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sync"
	"time"
)

var AppAsciiArt = fmt.Sprintf(""+
	"   _____       _ _           \n"+
	"  / ____|     (_|_)          \n"+
	" | |  __  ___  _ _ _ __ __ _ \n"+
	" | | |_ |/ _ \\| | | '__/ _` |\n"+
	" | |__| | (_) | | | | | (_| |\n"+
	"  \\_____|\\___/| |_|_|  \\__,_|\n"+
	"             _/ |     v%s \n"+
	"            |__/             \n\n", projectVersion)

var WorklogsCommand = &cli.Command{
	Name:  "worklogs",
	Usage: "Edit your today's work log",
	Action: func(c *cli.Context) error {
		newUi()
		loadWorklogs()
		err := app.ui.app.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func NewWorklogIssues() error {
	// goroutine awesomeness
	var err error
	startDate, endDate := MonthRange(app.time)
	if app.workLogsIssues.startDate == startDate && app.workLogsIssues.endDate == endDate {
		return nil
	}
	if app.workLogsIssues.startDate != startDate || app.workLogsIssues.endDate != endDate {
		app.ui.loaderView.Show("Fetching worklogs...")
		app.workLogs, err = GetWorklogs(MonthRange(app.time))
		app.ui.loaderView.Hide()
		if err != nil {
			return err
		}
		app.ui.calendar.update()
		app.ui.summary.update()
	}
	app.workLogsIssues.startDate = startDate
	app.workLogsIssues.endDate = endDate
	app.workLogsIssues.issues = []WorklogIssue{}
	waitGroup := sync.WaitGroup{}
	var errors []error
	errCh := make(chan error, len(app.workLogs.logs))
	for i := range app.workLogs.logs {
		waitGroup.Add(1)
		go func(workLog *Worklog) {
			issue, err := NewJiraClient().GetIssue(workLog.Issue.Key)
			if err != nil {
				errCh <- err // Send the error to the channel.
				return
			}
			app.workLogsIssues.issues = append(app.workLogsIssues.issues, WorklogIssue{Worklog: workLog, Issue: issue})
			waitGroup.Done()
		}(app.workLogs.logs[i])
	}
	waitGroup.Wait()
	close(errCh)
	// Collect all the errors.
	for err := range errCh {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

var IssuesCommand = &cli.Command{
	Name:  "issues",
	Usage: "Show currently assigned issues",
	Action: func(context *cli.Context) error {
		lastTickets, err := NewJiraClient().GetLatestIssues()
		if err != nil {
			return err
		}
		todaysWorklogs, err := GetWorklogs(DayRange(app.time))
		if err != nil {
			return err
		}
		alreadyLoggedIssues := []string{}
		for _, worklog := range todaysWorklogs.logs {
			alreadyLoggedIssues = append(alreadyLoggedIssues, worklog.Issue.Key)
		}
		todaysIssues, err := GetIssuesByKeys(alreadyLoggedIssues)
		if err != nil {
			return err
		}
		// FIXME - recent worklogs with those already logged
		issue, err := SelectIssueForm(append(todaysIssues.Issues, lastTickets.Issues...))
		if err != nil {
			return err
		}
		timeSpent, err := InputTimeSpentForm(issue)
		if err != nil {
			return err
		}
		err = spinner.New().Title("Logging work...").Action(func() {
			err = issue.LogWork(app.time, timeSpent)
		}).Run()
		if err != nil {
			return err
		}
		fmt.Printf("Successfully logged %s to ticket %s ", timeSpent, issue.Key)
		return nil
	},
}

var ViewIssueCommand = &cli.Command{
	Name:   "view",
	Usage:  "View issue in browser",
	Action: ViewIssueInBrowserAction,
}

var LogWorkCommand = &cli.Command{
	Name:      "log",
	Usage:     "Log work to specified issue",
	ArgsUsage: "ISSUE [TIME_SPENT]",
	Action: func(context *cli.Context) error {
		issueKey := ResolveIssueKey(context)
		timeSpent := context.Args().Get(1)
		if issueKey == "" {
			log.Fatalln("No issue key given / detected in git branch.")
		}
		issue, err := NewJiraClient().GetIssue(issueKey)
		if err != nil {
			return err
		}
		if timeSpent == "" {
			timeSpent, err = InputTimeSpentForm(issue)
			if err != nil {
				return err
			}
		}
		err = spinner.New().Title("Logging work...").Action(func() {
			err = issue.LogWork(app.time, timeSpent)
		}).Run()
		if err != nil {
			return err
		}
		fmt.Printf("Successfully logged %s to ticket %s ", timeSpent, issue.Key)
		return nil
	},
}

var DefaultAction = func(c *cli.Context) error {
	if c.Args().Get(0) != "" {
		fmt.Printf("Command not found: %v\n", c.Args().Get(0))
		os.Exit(1)
	}
	ticketFromBranch := ResolveIssueKey(c)
	if ticketFromBranch != "" {
		c.App.Metadata["JiraIssue"] = ticketFromBranch
		fmt.Printf("Detected possible ticket in git branch name - %s\n", ticketFromBranch)
		action, err := SelectActionForm([]string{"Log Work", "View Issue"})
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return nil
		}
		fmt.Printf("You choose %q\n", action)
		if action == "Log Work" {
			return GitOrIssueListAction(c) //fixme pass resolved Issue in context
		}
		if action == "View Issue" {
			return ViewIssueInBrowserAction(c) //fixme pass resolved Issue in context
		}
	}

	return GitOrIssueListAction(c)
}

var GitOrIssueListAction = func(c *cli.Context) error {
	issueKey := ResolveIssueKey(c)
	if issueKey != "" {
		issue, err := NewJiraClient().GetIssue(issueKey)
		if err != nil {
			return err
		}
		timeSpent, err := InputTimeSpentForm(issue)
		if err != nil {
			return nil
		}
		err = spinner.New().Title("Logging work...").Action(func() {
			err = issue.LogWork(app.time, timeSpent)
		}).Run()
		if err != nil {
			return err
		}
		fmt.Printf("Successfully logged %s to ticket %s ", timeSpent, issue.Key)
		return nil
	}

	err := IssuesCommand.Action(c)
	if err != nil {
		return err
	}
	return nil
}

var ViewIssueInBrowserAction = func(c *cli.Context) error {
	issueKey := ResolveIssueKey(c)
	if issueKey != "" {
		OpenURL(fmt.Sprintf("%s/browse/%s", Config.JiraUrl, issueKey))
	}
	return nil
}

var ConfigCommand = &cli.Command{
	Name:  "config",
	Usage: "configuration help",
	Action: func(context *cli.Context) error {
		//nolint:lll
		fmt.Print(`gojira needs a couple of env variables right now that you have to configure:
#1 Export below values in your .bashrc / .zshrc / .profile file:

export GOJIRA_JIRA_INSTANCE_URL="https://<INSTANCE>.atlassian.net"
export GOJIRA_JIRA_LOGIN="your@email.com"
export GOJIRA_JIRA_TOKEN= generate it at https://id.atlassian.com/manage-profile/security/api-tokens
export GOJIRA_TEMPO_TOKEN= generate it at https://<INSTANCE>.atlassian.net/plugins/servlet/ac/io.tempo.jira/tempo-app#!/configuration/api-integration

#2 Now we need to fetch one last env variable using previously saved values:
export GOJIRA_JIRA_ACCOUNT_ID= fetch it using this curl: 
curl --request GET \
  --url "$GOJIRA_JIRA_INSTANCE_URL/rest/api/3/user/bulk/migration?username=$GOJIRA_JIRA_LOGIN" \
  --header "Authorization: Basic $(echo -n $GOJIRA_JIRA_LOGIN:$GOJIRA_JIRA_TOKEN | base64)"

Save it and you should ready to go!
`)
		return nil
	},
}

func (issue Issue) LogWork(logTime *time.Time, timeSpent string) error {
	logrus.Infof("Logging %s of time to ticket %s at %s", timeSpent, issue.Key, logTime)
	todayWorklog, err := app.workLogs.LogsOnDate(logTime)
	if err != nil {
		return err
	}
	if Config.UpdateExistingWorklog {
		for index, workLog := range todayWorklog {
			if workLog.Issue.Key == issue.Key {
				//fmt.Println("Updating existing worklog...")
				timeSpentSum := FormatTimeSpent(TimeSpentToSeconds(timeSpent) + workLog.TimeSpentSeconds)
				err := todayWorklog[index].Update(timeSpentSum)
				if err != nil {
					return err
				}
				// FIXME - this should be only in cli mode
				//fmt.Printf("Successfully logged %s of time to ticket %s\n", timeSpent, workLog.Issue.Key)
				//fmt.Printf("Currently logged time: %s\n", FormatTimeSpent(CalculateTimeSpent(todayWorklog)))
				return nil
			}
		}
	}
	worklog, err := NewWorklog(issue.Key, logTime, timeSpent)
	if err != nil {
		return err
	}
	// add this workload to global object
	app.workLogs.logs = append(app.workLogs.logs, &worklog)
	app.workLogsIssues.issues = append(app.workLogsIssues.issues, WorklogIssue{Issue: issue, Worklog: &worklog})
	return nil
}
