package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type DayView struct {
	worklogList        *tview.Table
	worklogStatus      *tview.TextView
	latestIssuesList   *tview.Table
	latestIssuesStatus *tview.TextView
}

func NewDayView() *DayView {
	dayView := &DayView{
		worklogList: tview.NewTable(),
		worklogStatus: tview.NewTextView().SetChangedFunc(func() {
			app.ui.app.Draw()
		}),
		latestIssuesList: tview.NewTable(),
		latestIssuesStatus: tview.NewTextView().SetChangedFunc(func() {
			app.ui.app.Draw()
		}),
	}

	// FIXME instead border we could color code it or add some prompt to given section
	dayView.worklogList.SetBorder(true)
	dayView.latestIssuesList.SetBorder(true)
	dayView.latestIssuesList.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGray))
	dayView.worklogStatus.SetText(
		fmt.Sprintf("Worklogs - %s - [?h[white]]",
			app.time.Format("2006-01-02"),
		)).SetDynamicColors(true)

	flexView := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dayView.worklogStatus, 1, 1, false).
		AddItem(dayView.worklogList, 0, 1, true).
		AddItem(dayView.latestIssuesStatus, 1, 1, false).
		AddItem(dayView.latestIssuesList, 0, 1, false)

	dayView.worklogList.SetCell(0, 0, // FIXME use enums for column names
		tview.NewTableCell("Loading...").SetAlign(tview.AlignLeft),
	)

	// Make tab key able to switch between the two tables
	// Change focues table active row color to yellow and inactive to white
	flexView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// FIXME - it's not ideal - we should to check if given table is focused instead
		if event.Key() == tcell.KeyTab {
			if app.ui.app.GetFocus() == dayView.worklogList {
				app.ui.app.SetFocus(dayView.latestIssuesList)
				dayView.latestIssuesList.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
				dayView.worklogList.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite))
				return nil
			}
			app.ui.app.SetFocus(dayView.worklogList)
			dayView.worklogStatus.SetText(fmt.Sprintf(">%s", dayView.worklogStatus.GetText(true)))
			dayView.worklogList.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
			dayView.latestIssuesList.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite))
			return nil
		}
		return event
	})

	app.ui.pages.AddPage("worklog-view", flexView, true, true)

	dayView.worklogList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyRight:
			timePeriod := -time.Hour * 24
			if event.Key() == tcell.KeyRight {
				timePeriod = time.Hour * 24
			}
			newTime := app.time.Add(timePeriod)
			app.time = &newTime
			loadWorklogs()
			app.ui.calendar.update()
			break
		}
		return event
	})

	dayView.loadLatest()

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
	d.worklogList.Clear()
	d.worklogList.SetSelectable(true, false)
	color := tcell.ColorWhite
	for r := 0; r < len(logs); r++ {
		d.worklogList.SetCell(r, 0, // FIXME use enums for column names
			tview.NewTableCell((logs)[r].Issue.Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.worklogList.SetCell(r, 1,
			tview.NewTableCell((logs)[r].Issue.Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.worklogList.SetCell(r, 2,
			tview.NewTableCell(FormatTimeSpent((logs)[r].WorkLog.TimeSpentSeconds)).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	d.worklogList.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		NewUpdateWorklogForm(d, logs, row)
	})
	timeSpent := CalculateTimeSpent(getWorkLogsFromWorkLogIssues(logs))
	d.worklogStatus.SetText(
		fmt.Sprintf("Worklogs - %s - [%s%s[white]]",
			app.time.Format("2006-01-02"),
			GetTimeSpentColorTag(timeSpent),
			FormatTimeSpent(timeSpent),
		)).SetDynamicColors(true)
}

func (d *DayView) loadLatest() {
	d.latestIssuesStatus.SetText("Latest issues").SetDynamicColors(true)
	issues, err := GetLatestIssues()
	if err != nil {
		app.ui.errorView.ShowError(err.Error())
		return
	}
	d.latestIssuesList.Clear()
	d.latestIssuesList.SetSelectable(true, false)
	color := tcell.ColorWhite
	for r := 0; r < len(issues.Issues); r++ {
		d.latestIssuesList.SetCell(r, 0, // FIXME use enums for column names
			tview.NewTableCell((issues.Issues)[r].Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.latestIssuesList.SetCell(r, 1,
			tview.NewTableCell((issues.Issues)[r].Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	d.latestIssuesList.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		NewAddWorklogForm(d, issues.Issues, row)
	})
}

func NewAddWorklogForm(d *DayView, issues []Issue, row int) *tview.Form {
	var form *tview.Form

	newWorklog := func() {
		timeSpent := form.GetFormItem(0).(*tview.InputField).GetText()
		app.ui.flex.SetTitle(" gojira - adding worklog... ")
		go func() {
			issue, err := GetIssue(issues[row].Key)
			if err != nil {
				app.ui.errorView.ShowError(err.Error())
				return
			}
			err = issue.LogWork(timeSpent)
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
		AddInputField("Time spent", "", 20, nil, nil).
		AddButton("Add", newWorklog).
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
	form.SetBorder(true).SetTitle("New worklog").SetTitleAlign(tview.AlignLeft)
	app.ui.pages.AddPage("worklog-form", form, true, true)
	return form
}

func NewUpdateWorklogForm(d *DayView, workLogIssues []*WorkLogIssue, row int) *tview.Form {
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
