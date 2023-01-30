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
		var workLogIssues []WorkLogIssue
		// goroutine awesomeness
		waitGroup := sync.WaitGroup{}
		for _, workLog := range GetWorkLogs() {
			waitGroup.Add(1)
			go func(workLog WorkLog) {
				workLogIssues = append(workLogIssues, WorkLogIssue{WorkLog: workLog, Issue: GetIssue(workLog.Issue.Key)})
				waitGroup.Done()
			}(workLog)
		}
		waitGroup.Wait()

		if len(workLogIssues) == 0 {
			fmt.Println("You don't have any logged work today.")
			return nil
		}
		newUi()
		newWorkLogTable(workLogIssues)
		err := app.ui.app.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

var IssuesCommand = &cli.Command{
	Name:  "issues",
	Usage: "Show currently assigned issues",
	Action: func(context *cli.Context) error {
		lastTickets := GetLatestIssues()
		issue, err := PromptForIssueSelection(lastTickets.Issues)
		if err != nil {
			return nil
		}
		timeSpent, err := PromptForTimeSpent("Add work log")
		if err != nil {
			return nil
		}
		issue.LogWork(timeSpent)
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
		issue := GetIssue(issueKey)
		fmt.Printf("%s %s\n", issue.Key, issue.Fields.Summary)
		fmt.Printf("Status: %s\n", issue.Fields.Status.Name)
		if timeSpent == "" {
			var err error
			timeSpent, err = PromptForTimeSpent("Add work log")
			if err != nil {
				log.Fatalln(err)
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
		issue := GetIssue(issueKey)
		fmt.Printf("Status: %s\nSummary: %s\n", issue.Fields.Status.Name, issue.Fields.Summary)
		// log time or view issue
		timeSpent, err := PromptForTimeSpent("Add work log")
		if err != nil {
			return nil
		}
		issue.LogWork(timeSpent)
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

func (issue Issue) LogWork(timeSpent string) {
	workLogs := GetWorkLogs()
	if Config.UpdateExistingWorkLog {
		for index, workLog := range workLogs {
			if workLog.Issue.Key == issue.Key {
				fmt.Println("Updating existing worklog...")
				timeSpentSum := FormatTimeSpent(TimeSpentToSeconds(timeSpent) + workLog.TimeSpentSeconds)
				workLogs[index].Update(timeSpentSum)
				fmt.Printf("Successfully logged %s of time to ticket %s\n", timeSpent, workLog.Issue.Key)
				fmt.Printf("Currently logged time: %s\n", CalculateTimeSpent(workLogs))
				return
			}
		}
	}
	issue.NewWorkLog(timeSpent)
	// naive issue struct for quicker summary
	workLogs = append(workLogs, WorkLog{TimeSpentSeconds: TimeSpentToSeconds(timeSpent)})
	fmt.Printf("Currently logged time: %s\n", CalculateTimeSpent(workLogs))
}
