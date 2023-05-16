package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type UserInteface struct {
	app      *tview.Application
	flex     *tview.Flex
	pages    *tview.Pages
	table    *tview.Table
	status   *tview.TextView
	modal    *tview.Modal
	calendar *Calendar
}

var running = make(chan bool, 1)

func newUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()
	app.ui.table = tview.NewTable()
	app.ui.calendar = NewCalendar()
	app.ui.status = tview.NewTextView().SetChangedFunc(func() {
		app.ui.app.Draw()
	})
	app.ui.modal = tview.NewModal().SetText("Something went wrong")
	app.ui.modal.SetTitle("Error!")
	app.ui.modal.AddButtons([]string{"OK"})
	app.ui.modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			app.ui.pages.RemovePage("error")
		}
		return event
	})

	app.ui.flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(app.ui.status, 1, 1, false).
			AddItem(app.ui.pages, 0, 1, true),
			80, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(app.ui.calendar.Table, 28, 1, false),
			30, 1, false)
	app.ui.flex.SetBorder(true).SetTitle("gojira")
	app.ui.app.SetRoot(app.ui.flex, true)

	app.ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'n', 'p':
				timePeriod := -time.Hour * 24
				if event.Rune() == 'n' {
					timePeriod = time.Hour * 24
				}
				app.time = app.time.Add(timePeriod)
				app.ui.table.Clear()
				app.ui.table.SetCell(0, 0, tview.NewTableCell("Loading..."))
				select {
				case running <- true:
					go func() {
						defer func() { <-running }()
						GetWorkLogIssues()
						logs, _ := app.workLogsIssues.IssuesOnDate(app.time)
						newWorkLogView(logs)
					}()
				default:
					// The goroutine is already running, do nothing
				}
				app.ui.calendar.update()
				break
			}
		}
		return event
	})
}

func newWorkLogView(workLogs []*WorkLogIssue) {
	app.ui.pages.AddPage("worklog-view", app.ui.table, true, true)
	app.ui.table.Clear()
	app.ui.table.SetSelectable(true, false)
	color := tcell.ColorWhite
	for r := 0; r < len(workLogs); r++ {
		app.ui.table.SetCell(r, 0, // FIXME use enums for column names
			tview.NewTableCell((workLogs)[r].Issue.Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		app.ui.table.SetCell(r, 1,
			tview.NewTableCell((workLogs)[r].Issue.Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		app.ui.table.SetCell(r, 2,
			tview.NewTableCell(FormatTimeSpent((workLogs)[r].WorkLog.TimeSpentSeconds)).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	app.ui.table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		newWorklogForm(workLogs, row)
	})
	timeSpent := CalculateTimeSpent(getWorkLogsFromWorkLogIssues(workLogs))
	app.ui.status.SetText(
		fmt.Sprintf("Worklogs - %s -  [%s%s[white]]",
			app.time.Format("2006-01-02"),
			GetTimeSpentColorTag(timeSpent),
			FormatTimeSpent(timeSpent),
		)).SetDynamicColors(true)
	app.ui.pages.ShowPage("worklog-view")
}

func newWorklogForm(workLogIssues []*WorkLogIssue, row int) *tview.Form {
	var form *tview.Form

	updateWorklog := func() {
		timeSpent := form.GetFormItem(0).(*tview.InputField).GetText()
		workLogIssues[row].WorkLog.Update(timeSpent)
		app.ui.pages.HidePage("worklog-form")
		newWorkLogView(workLogIssues)
	}

	form = tview.NewForm().
		AddInputField("Time spent", FormatTimeSpent(workLogIssues[row].WorkLog.TimeSpentSeconds), 20, nil, nil).
		AddButton("Update", updateWorklog).
		AddButton("Cancel", func() { // FIXME can't move to cancel button
			app.ui.pages.HidePage("worklog-form")
		})
	form.SetBorder(true).SetTitle("Update worklog").SetTitleAlign(tview.AlignLeft)
	app.ui.pages.AddPage("worklog-form", form, true, true)
	return form
}
