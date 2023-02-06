package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type UserInteface struct {
	app    *tview.Application
	flex   *tview.Flex
	pages  *tview.Pages
	table  *tview.Table
	status *tview.TextView
}

func newUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()
	app.ui.table = tview.NewTable()
	app.ui.status = tview.NewTextView().SetText(
		fmt.Sprintf("worklogs - %s", app.time.Format("2006-01-02")),
	).SetChangedFunc(func() {
		app.ui.app.Draw()
	})
	app.ui.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.ui.status, 1, 1, false).
		AddItem(app.ui.pages, 0, 1, true)
	app.ui.flex.SetBorder(true).SetTitle("gojira")

	app.ui.app.SetRoot(app.ui.flex, true)
	app.ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'p':
				app.time = app.time.Add(-time.Hour * 24)
				newWorkLogView(GetWorkLogIssues())
			}
		}
		return event
	})
}

func newWorkLogView(workLogs []WorkLogIssue) {
	app.ui.pages.AddPage("worklog-view", app.ui.table, true, true)
	app.ui.table.Clear()
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
	app.ui.status.SetText(fmt.Sprintf("worklogs - %s", app.time.Format("2006-01-02")))
	app.ui.pages.ShowPage("worklog-view")
}

func newWorklogForm(workLogIssues []WorkLogIssue, row int) *tview.Form {
	var form *tview.Form

	updateWorklog := func() {
		timeSpent := form.GetFormItem(0).(*tview.InputField).GetText()
		app.ui.status.SetText("Updating....")
		workLogIssues[row].WorkLog.Update(timeSpent)
		app.ui.pages.HidePage("worklog-form")
		newWorkLogView(workLogIssues)
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
