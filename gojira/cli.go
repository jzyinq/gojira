package gojira

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sync"
)

var WorkLogsCommand = &cli.Command{
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

func NewWorkLogIssues() error {
	// goroutine awesomeness
	var err error
	startDate, endDate := MonthRange(app.time)
	if app.workLogsIssues.startDate == startDate && app.workLogsIssues.endDate == endDate {
		return nil
	}
	if app.workLogsIssues.startDate != startDate || app.workLogsIssues.endDate != endDate {
		app.workLogs, err = GetWorkLogs()
		if err != nil {
			return err
		}
		app.ui.calendar.update()
		app.ui.summary.update()
	}
	app.workLogsIssues.startDate = startDate
	app.workLogsIssues.endDate = endDate
	app.workLogsIssues.issues = []WorkLogIssue{}
	waitGroup := sync.WaitGroup{}
	var errors []error
	errCh := make(chan error, len(app.workLogs.logs))
	for i, _ := range app.workLogs.logs {
		waitGroup.Add(1)
		go func(workLog *WorkLog) {
			issue, err := GetIssue(workLog.Issue.Key)
			if err != nil {
				errCh <- err // Send the error to the channel.
				return
			}
			app.workLogsIssues.issues = append(app.workLogsIssues.issues, WorkLogIssue{WorkLog: workLog, Issue: issue})
			waitGroup.Done()
		}(&app.workLogs.logs[i])
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
		lastTickets, err := GetLatestIssues()
		if err != nil {
			return err
		}
		issue, err := PromptForIssueSelection(lastTickets.Issues)
		if err != nil {
			return err
		}
		timeSpent, err := PromptForTimeSpent("Add work log")
		if err != nil {
			return err
		}
		err = issue.LogWork(timeSpent)
		if err != nil {
			return err
		}
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
		issue, err := GetIssue(issueKey)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s\n", issue.Key, issue.Fields.Summary)
		fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
		if timeSpent == "" {
			timeSpent, err = PromptForTimeSpent("Add work log")
			if err != nil {
				return err
			}
		}

		issue.LogWork(timeSpent)
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
		prompt := promptui.Select{
			Label: "Select Action",
			Items: []string{"Log Work", "View Issue"},
		}
		_, action, err := prompt.Run()

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
		issue, err := GetIssue(issueKey)
		if err != nil {
			return err
		}
		fmt.Printf("Status: %s\nSummary: %s\n", issue.Fields.Status.Name, issue.Fields.Summary)
		// log time or view issue
		timeSpent, err := PromptForTimeSpent("Add work log")
		if err != nil {
			return nil
		}
		err = issue.LogWork(timeSpent)
		if err != nil {
			return err
		}
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
		OpenUrl(fmt.Sprintf("%s/browse/%s", Config.JiraUrl, issueKey))
	}
	return nil
}

var ConfigCommand = &cli.Command{
	Name:  "config",
	Usage: "configuration help",
	Action: func(context *cli.Context) error {
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

func (issue Issue) LogWork(timeSpent string) error {
	todayWorklog, _ := app.workLogs.LogsOnDate(app.time) // FIXME error handling
	if Config.UpdateExistingWorkLog {
		for index, workLog := range todayWorklog {
			if workLog.Issue.Key == issue.Key {
				fmt.Println("Updating existing worklog...")
				timeSpentSum := FormatTimeSpent(TimeSpentToSeconds(timeSpent) + workLog.TimeSpentSeconds)
				err := todayWorklog[index].Update(timeSpentSum)
				if err != nil {
					return err
				}
				fmt.Printf("Successfully logged %s of time to ticket %s\n", timeSpent, workLog.Issue.Key)
				fmt.Printf("Currently logged time: %s\n", FormatTimeSpent(CalculateTimeSpent(todayWorklog)))
				return nil
			}
		}
	}
	err := issue.NewWorkLog(timeSpent)
	if err != nil {
		return err
	}
	// naive issue struct for quicker summary
	loggedWorklog := &WorkLog{TimeSpentSeconds: TimeSpentToSeconds(timeSpent)}
	todayWorklog = append(todayWorklog, loggedWorklog)
	// append new WorkLogIssue to app.workLogsIssues.issues with newly recorded worklog
	// FIXME - newly added worklog is not showing up correcly - updated one does right
	app.workLogsIssues.issues = append(app.workLogsIssues.issues, WorkLogIssue{Issue: issue, WorkLog: loggedWorklog})

	fmt.Printf("Currently logged time: %s\n", FormatTimeSpent(CalculateTimeSpent(todayWorklog)))
	return nil
}
