package gojira

import (
	"fmt"
	"github.com/rivo/tview"
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
	summary.SetText("?h/?h")
	summary.SetTextAlign(tview.AlignRight)
	return summary
}

func (s *Summary) update() {
	s.SetText(
		fmt.Sprintf("%s/%dh",
			FormatTimeSpent(app.workLogs.TotalTimeSpent()), workingHoursInMonth(app.time.Year(), app.time.Month()),
		),
	)
}
