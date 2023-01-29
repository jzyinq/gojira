package gojira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UserInteface struct {
	app   *tview.Application
	frame *tview.Frame
	pages *tview.Pages
	table *tview.Table
}

func newUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()
	app.ui.table = tview.NewTable()
	app.ui.pages.AddPage("worklog-view", app.ui.table, true, true)
	app.ui.frame = tview.NewFrame(app.ui.pages)
	app.ui.frame.SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Worklogs", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("gojira v0.0.9", true, tview.AlignRight, tcell.ColorWhite).
		AddText("(p)revious day   (n)ext day", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("(d)elete worklog (enter) update worklog", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("Status...", false, tview.AlignCenter, tcell.ColorGreen)
	app.ui.app.SetRoot(app.ui.frame, true)

	// TODO worklog view

}

func newWorkLogTable(workLogs []WorkLogIssue) {
	app.ui.table.SetSelectable(true, false)
	color := tcell.ColorWhite
	for r := 0; r < len(workLogs); r++ {
		app.ui.table.SetCell(r, 0, // FIXME use enums for column names
			tview.NewTableCell(workLogs[r].Issue.Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		app.ui.table.SetCell(r, 1,
			tview.NewTableCell(workLogs[r].Issue.Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		app.ui.table.SetCell(r, 2,
			tview.NewTableCell(FormatTimeSpent(workLogs[r].WorkLog.TimeSpentSeconds)).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	app.ui.table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		newWorklogForm(workLogs, row)
	})
	app.ui.pages.ShowPage("worklog-view")
}

func newWorklogForm(workLogIssues []WorkLogIssue, row int) *tview.Form {
	var form *tview.Form

	updateWorklog := func() {
		timeSpent := form.GetFormItem(0).(*tview.InputField).GetText()
		workLogIssues[row].WorkLog.Update(timeSpent)
		app.ui.pages.HidePage("worklog-form")
		newWorkLogTable(workLogIssues)
	}

	form = tview.NewForm().
		AddInputField("Time spent", FormatTimeSpent(workLogIssues[row].WorkLog.TimeSpentSeconds), 20, nil, nil).
		AddButton("Update", updateWorklog).
		AddButton("Cancel", func() {
			app.ui.pages.HidePage("worklog-form")
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	app.ui.pages.AddPage("worklog-form", form, true, true)
	return form
}
