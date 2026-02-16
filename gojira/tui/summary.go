package gojira

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type Summary struct {
	*tview.TextView
}

func NewSummary() *Summary {
	summary := &Summary{
		TextView: tview.NewTextView().SetChangedFunc(func() {
			App.ui.app.Draw()
		}),
	}
	summary.SetText("Loading...")
	summary.SetTextAlign(tview.AlignCenter)
	return summary
}

func (s *Summary) update() {
	totalTimeSpent := FormatTimeSpent(App.workLogs.TotalTimeSpentToPresentDay())
	// that's a hack to remove spaces between hours and minutes
	totalTimeSpent = strings.Join(strings.Fields(totalTimeSpent), "")
	workingHours := workingHoursInMonthToPresentDay(App.time.Year(), App.time.Month())
	difference := workingHoursAbsoluteDiff(workingHours)
	status := fmt.Sprintf("Total %s/%dh", totalTimeSpent, workingHours)
	if difference != 0 {
		status = fmt.Sprintf("Total %s/%dh (%s)", totalTimeSpent, workingHours, FormatTimeSpent(difference))
	}
	s.SetText(status)
	s.SetTextColor(GetTimeSpentColor(App.workLogs.TotalTimeSpentToPresentDay(), workingHours))
}

func workingHoursInMonthToPresentDay(year int, month time.Month) int {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	totalWorkHours := 0

	for t.Month() == month && t.Before(time.Now().Local()) {
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday && !app.holidays.IsHoliday(&t) {
			totalWorkHours += 8
		}
		t = t.AddDate(0, 0, 1)
	}
	return totalWorkHours
}

func workingHoursAbsoluteDiff(workingHours int) int {
	difference := workingHours*60*60 - App.workLogs.TotalTimeSpentToPresentDay()
	if difference < 0 {
		difference = -difference
	}
	return difference
}
