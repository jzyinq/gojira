package gojira

import (
	"fmt"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type Summary struct {
	*tview.TextView
}

func NewSummary() *Summary {
	summary := &Summary{
		TextView: tview.NewTextView().SetChangedFunc(func() {
			app.ui.app.Draw()
		}),
	}
	summary.SetText("Loading...")
	summary.SetTextAlign(tview.AlignCenter)
	return summary
}

func (s *Summary) update() {
	totalTimeSpent := FormatTimeSpent(app.workLogs.TotalTimeSpentToPresentDay())
	// that's a hack to remove spaces between hours and minutes
	totalTimeSpent = strings.Join(strings.Fields(totalTimeSpent), "")
	workingHours := workingHoursInMonthToPresentDay(app.time.Year(), app.time.Month())
	difference := workingHoursAbsoluteDiff(workingHours)
	status := fmt.Sprintf("Total %s/%dh", totalTimeSpent, workingHours)
	if difference != 0 {
		status = fmt.Sprintf("Total %s/%dh (%s)", totalTimeSpent, workingHours, FormatTimeSpent(difference))
	}
	s.SetText(status)
	s.SetTextColor(GetTimeSpentColor(app.workLogs.TotalTimeSpentToPresentDay(), workingHours))
}

func workingHoursInMonthToPresentDay(year int, month time.Month) int {
	holidays := NewHolidays("PL")
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	totalWorkHours := 0

	for t.Month() == month && t.Before(time.Now().Local()) {
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday && !holidays.IsHoliday(&t) {
			totalWorkHours += 8
		}
		t = t.AddDate(0, 0, 1)
	}
	return totalWorkHours
}

func workingHoursAbsoluteDiff(workingHours int) int {
	difference := workingHours*60*60 - app.workLogs.TotalTimeSpentToPresentDay()
	if difference < 0 {
		difference = -difference
	}
	return difference
}
