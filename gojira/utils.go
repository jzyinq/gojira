package gojira

import (
	"fmt"
	"math"
	"os/exec"
	"regexp"
)

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
		timeSpent = timeSpent + fmt.Sprintf("%vm", math.Round(floatPart*60))
	}
	return timeSpent
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
