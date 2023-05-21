package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type DayView struct {
	table  *tview.Table
	status *tview.TextView
}

func NewDayView() *DayView {
	dayView := &DayView{
		table: tview.NewTable(),
		status: tview.NewTextView().SetChangedFunc(func() {
			app.ui.app.Draw()
		}),
	}

	app.ui.pages.AddPage("worklog-view",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(dayView.status, 1, 1, false).
			AddItem(dayView.table, 0, 1, true),
		true, true)

	dayView.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyRight:
			timePeriod := -time.Hour * 24
			if event.Key() == tcell.KeyRight {
				timePeriod = time.Hour * 24
			}
			app.time = app.time.Add(timePeriod)
			loadWorklogs()
			app.ui.calendar.update()
			break
		}
		return event
	})

	return dayView
}

var loadingWorklogs = make(chan bool, 1)

func loadWorklogs() {
	select {
	case loadingWorklogs <- true:
		go func() {
			defer func() { <-loadingWorklogs }()
			app.ui.flex.SetTitle(" gojira - fetching data... ")
			err := NewWorkLogIssues()
			if err != nil {
				app.ui.errorView.ShowError(err.Error())
			}
			app.ui.dayView.update()
			app.ui.flex.SetTitle(" gojira ")
		}()
	default:
		// The goroutine is already loadingWorklogs, do nothing
	}
}

func (d *DayView) update() {
	logs, _ := app.workLogsIssues.IssuesOnDate(app.time)
	d.table.Clear()
	d.table.SetSelectable(true, false)
	color := tcell.ColorWhite
	for r := 0; r < len(logs); r++ {
		d.table.SetCell(r, 0, // FIXME use enums for column names
			tview.NewTableCell((logs)[r].Issue.Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.table.SetCell(r, 1,
			tview.NewTableCell((logs)[r].Issue.Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.table.SetCell(r, 2,
			tview.NewTableCell(FormatTimeSpent((logs)[r].WorkLog.TimeSpentSeconds)).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	d.table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		newWorklogForm(d, logs, row)
	})
	timeSpent := CalculateTimeSpent(getWorkLogsFromWorkLogIssues(logs))
	d.status.SetText(
		fmt.Sprintf("Worklogs - %s -  [%s%s[white]]",
			app.time.Format("2006-01-02"),
			GetTimeSpentColorTag(timeSpent),
			FormatTimeSpent(timeSpent),
		)).SetDynamicColors(true)
}

func newWorklogForm(d *DayView, workLogIssues []*WorkLogIssue, row int) *tview.Form {
	var form *tview.Form

	updateWorklog := func() {
		timeSpent := form.GetFormItem(0).(*tview.InputField).GetText()
		app.ui.flex.SetTitle(" gojira - updating worklog... ")
		go func() {
			err := workLogIssues[row].WorkLog.Update(timeSpent)
			app.ui.flex.SetTitle(" gojira ")
			if err != nil {
				app.ui.errorView.ShowError(err.Error())
				return
			}
			d.update()
			app.ui.pages.RemovePage("worklog-form")
			app.ui.calendar.update()
		}()
	}

	deleteWorklog := func() {
		app.ui.flex.SetTitle(" gojira - deleting worklog... ")
		go func() {
			err := app.workLogs.Delete(workLogIssues[row].WorkLog)
			app.ui.flex.SetTitle(" gojira ")
			if err != nil {
				app.ui.errorView.ShowError(err.Error())
				return
			}
			d.update()
			app.ui.pages.RemovePage("worklog-form")
			app.ui.calendar.update()
		}()
	}

	form = tview.NewForm().
		AddInputField("Time spent", FormatTimeSpent(workLogIssues[row].WorkLog.TimeSpentSeconds), 20, nil, nil).
		AddButton("Update", updateWorklog).
		AddButton("Delete", deleteWorklog).
		AddButton("Cancel", func() {
			app.ui.pages.RemovePage("worklog-form")
		})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.ui.pages.RemovePage("worklog-form")
			break
		}
		return event
	})
	form.SetBorder(true).SetTitle("Update worklog").SetTitleAlign(tview.AlignLeft)
	app.ui.pages.AddPage("worklog-form", form, true, true)
	return form
}
