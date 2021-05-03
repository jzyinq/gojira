package gojira

import (
	"fmt"
	"math"
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
