package gojira

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"math"
	"os/exec"
	"regexp"
)

func getWorkLogsFromWorkLogIssues(workLogIssues []WorkLogIssue) []WorkLog {
	var workLogs []WorkLog
	for _, workLog := range workLogIssues {
		workLogs = append(workLogs, workLog.WorkLog)
	}
	return workLogs
}

func CalculateTimeSpent(workLogs []WorkLog) string {
	timeSpentInSeconds := 0
	for _, workLog := range workLogs {
		timeSpentInSeconds += workLog.TimeSpentSeconds
	}
	return FormatTimeSpent(timeSpentInSeconds)
}

func FormatTimeSpent(timeSpentSeconds int) string {
	timeInHours := float64(timeSpentSeconds) / 60 / 60
	intPart, floatPart := math.Modf(timeInHours)
	timeSpent := ""
	if intPart > 0 {
		timeSpent = fmt.Sprintf("%vh", intPart)
	}
	if floatPart > 0 {
		if (intPart) > 0 {
			timeSpent += " "
		}
		timeSpent = timeSpent + fmt.Sprintf("%vm", math.Round(floatPart*60))
	}
	return timeSpent
}

func ResolveIssueKey(c *cli.Context) string {
	issueKey := ""
	if 	c.App.Metadata["JiraIssue"] != nil {
		issueKey = fmt.Sprintf("%s", c.App.Metadata["JiraTicket"])
	}
	issueKey = FindIssueKeyInString(c.Args().Get(0))
	if issueKey == "" {
		issueKey = GetTicketFromGitBranch()
	}

	return issueKey
}

func GetTicketFromGitBranch() string {
	gitBranch, err := exec.Command("git", "branch", "--show-current").CombinedOutput()
	if err != nil {
		return ""
	}
	return FindIssueKeyInString(string(gitBranch))
}

func FindIssueKeyInString(possibleUrl string) string {
	r, _ := regexp.Compile("([A-Z]+-[0-9]+)")
	match := r.FindString(possibleUrl)
	return match
}

func OpenUrl(url string) {
	// silence browser logs
	browser.Stdout = nil
	browser.Stderr = nil
	err := browser.OpenURL(url)
	if err != nil {
		fmt.Println(err)
		return
	}
}
