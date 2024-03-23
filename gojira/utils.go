package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
	"math"
	"os/exec"
	"regexp"
	"time"
)

const dateLayout = "2006-01-02"

func getWorkLogsFromWorkLogIssues(workLogIssues []*WorkLogIssue) []*WorkLog {
	var workLogs []*WorkLog
	for _, workLog := range workLogIssues {
		workLogs = append(workLogs, workLog.WorkLog)
	}
	return workLogs
}

func CalculateTimeSpent(workLogs []*WorkLog) int {
	timeSpentInSeconds := 0
	for _, workLog := range workLogs {
		timeSpentInSeconds += workLog.TimeSpentSeconds
	}
	return timeSpentInSeconds
}

func GetTimeSpentColorTag(timeSpentInSeconds int, hours int) string {
	switch {
	case timeSpentInSeconds < hours*60*60 && timeSpentInSeconds > 0:
		return "[orange]"
	case timeSpentInSeconds == hours*60*60:
		return "[green]"
	case timeSpentInSeconds > hours*60*60:
		return "[blue]"
	default:
		return "[white]"
	}
}

func GetTimeSpentColor(timeSpentInSeconds int, hours int) tcell.Color {
	switch {
	case timeSpentInSeconds < hours*60*60 && timeSpentInSeconds > 0:
		return tcell.ColorOrange
	case timeSpentInSeconds == hours*60*60:
		return tcell.ColorGreen
	case timeSpentInSeconds > hours*60*60:
		return tcell.ColorBlue
	default:
		return tcell.ColorWhite
	}
}

func FormatTimeSpent(timeSpentSeconds int) string {
	timeInHours := float64(timeSpentSeconds) / 60 / 60
	intPart, floatPart := math.Modf(timeInHours)
	timeSpent := ""
	if timeSpentSeconds == 0 {
		return "0"
	}
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
	//if c.App.Metadata["JiraIssue"] != nil {
	//	issueKey = fmt.Sprintf("%s", c.App.Metadata["JiraTicket"])
	//}
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

func FindIssueKeyInString(possibleURL string) string {
	r, _ := regexp.Compile("([A-Z]+-[0-9]+)")
	match := r.FindString(possibleURL)
	return match
}

func OpenURL(url string) {
	// silence browser logs
	browser.Stdout = nil
	browser.Stderr = nil
	err := browser.OpenURL(url)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func WeekRange(today time.Time) (time.Time, time.Time) {
	y, w := today.ISOWeek()
	firstDay := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
	for firstDay.Weekday() != time.Monday {
		firstDay = firstDay.AddDate(0, 0, -1)
	}

	for {
		y1, w1 := firstDay.ISOWeek()
		if y1 == y && w1 == w {
			break
		}
		firstDay = firstDay.AddDate(0, 0, 1)
	}

	lastDay := firstDay.AddDate(0, 0, 6) // Adding 6 days to get to Sunday
	return firstDay.Truncate(24 * time.Hour), lastDay.Truncate(24 * time.Hour)
}

func MonthRange(t *time.Time) (time.Time, time.Time) {
	firstDayOfCurrentMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	firstDayOfNextMonth := firstDayOfCurrentMonth.AddDate(0, 1, 0)
	return firstDayOfCurrentMonth, firstDayOfNextMonth
}

func workingHoursInMonthToPresentDay(year int, month time.Month) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	totalWorkHours := 0

	for t.Month() == month && t.Before(time.Now().Local()) {
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
			totalWorkHours += 8
		}
		t = t.AddDate(0, 0, 1)
	}
	return totalWorkHours
}
