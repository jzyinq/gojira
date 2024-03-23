package gojira

import (
	"fmt"
	"github.com/rivo/tview"
	"strings"
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
	summary.SetText("Calendar ?h/?h")
	summary.SetTextAlign(tview.AlignCenter)
	return summary
}

func (s *Summary) update() {
	totalTimeSpent := FormatTimeSpent(app.workLogs.TotalTimeSpent())
	// that's a hack to remove spaces between hours and minutes
	totalTimeSpent = strings.Join(strings.Fields(totalTimeSpent), "")
	workingHours := workingHoursInMonth(app.time.Year(), app.time.Month())
	s.SetText(fmt.Sprintf("Monthly %s/%dh", totalTimeSpent, workingHours))
}
