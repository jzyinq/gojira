package gojira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UserInteface struct {
	app   *tview.Application
	frame *tview.Frame
	pages *tview.Pages
}

func NewUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()
	app.ui.frame = tview.NewFrame(app.ui.pages)
	app.ui.app.SetRoot(app.ui.frame, true)
}

func NewWorkLogView(workLogs []WorkLogIssue) {
	table := tview.NewTable().SetSelectable(true, false)
	color := tcell.ColorWhite
	// FIXME  Set fixed number of rows
	for r := 0; r < len(workLogs); r++ {
		table.SetCell(r, 0, // FIXME use enums for column names
			tview.NewTableCell(workLogs[r].Issue.Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		table.SetCell(r, 1,
			tview.NewTableCell(workLogs[r].Issue.Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		table.SetCell(r, 2,
			tview.NewTableCell(FormatTimeSpent(workLogs[r].WorkLog.TimeSpentSeconds)).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		// akcja na edycje tabeli
	})
	app.ui.pages.AddPage("worklog-view", table, true, true)
	app.ui.frame.SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Worklogs", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("gojira v0.0.9", true, tview.AlignRight, tcell.ColorWhite).
		AddText("(p)revious day   (n)ext day", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("(d)elete worklog (enter) update worklog", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("Status...", false, tview.AlignCenter, tcell.ColorGreen)
}
