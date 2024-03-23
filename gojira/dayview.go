package gojira

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"strings"
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
			// FIXME not necessary since you can't jump between calendar days on latest issues
			//dayView.worklogStatus.SetText(fmt.Sprintf("%s", dayView.worklogStatus.GetText(true)))
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
			logrus.Debug("Changing date to ", newTime)
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
			app.ui.loaderView.Show("Fetching worklogs...")
			err := NewWorkLogIssues()
			if err != nil {
				app.ui.errorView.ShowError(err.Error())
			}
			app.ui.dayView.update()
			app.ui.loaderView.Hide()
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
			GetTimeSpentColorTag(timeSpent, 8),
			FormatTimeSpent(timeSpent),
		)).SetDynamicColors(true)
}

func (d *DayView) loadLatest() {
	d.latestIssuesStatus.SetText("Latest issues").SetDynamicColors(true)
	issues, err := NewJiraClient().GetLatestIssues()
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

// DateRange is a struct for holding the start and end dates
type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

func ParseDateRange(dateStr string) (DateRange, error) {
	// Define the expected date layout
	const layout = "2006-01-02"
	var startDate, endDate time.Time
	var err error

	// Check if the string contains "->", indicating a range
	if strings.Contains(dateStr, "->") {
		dateParts := strings.Split(dateStr, "->")
		if len(dateParts) != 2 {
			return DateRange{}, errors.New("invalid date range format")
		}

		// Parse the start and end dates
		startDate, err = time.Parse(layout, dateParts[0])
		if err != nil {
			return DateRange{}, fmt.Errorf("error parsing start date: %w", err)
		}

		endDate, err = time.Parse(layout, dateParts[1])
		if err != nil {
			return DateRange{}, fmt.Errorf("error parsing end date: %w", err)
		}

	} else {
		// Parse the single date
		startDate, err = time.Parse(layout, dateStr)
		if err != nil {
			return DateRange{}, fmt.Errorf("error parsing date: %w", err)
		}
		endDate = startDate
	}

	return DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func NewAddWorklogForm(d *DayView, issues []Issue, row int) *tview.Form {
	var form *tview.Form

	newWorklog := func() {

		logTime := form.GetFormItem(0).(*tview.InputField).GetText()
		timeSpent := form.GetFormItem(1).(*tview.InputField).GetText()
		go func() {
			app.ui.loaderView.Show("Adding worklog...")
			defer app.ui.loaderView.Hide()
			issue, err := NewJiraClient().GetIssue(issues[row].Key)
			if err != nil {
				app.ui.errorView.ShowError(err.Error())
				return
			}
			// TODO use ParseDateRange and LogWork for each day in range
			dateRange, err := ParseDateRange(logTime)
			if err != nil {
				app.ui.errorView.ShowError(err.Error())
				return
			}
			for day := dateRange.StartDate; day.Before(dateRange.EndDate.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
				err := issue.LogWork(&day, timeSpent)
				if err != nil {
					app.ui.errorView.ShowError(err.Error())
					return
				}
			}
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
		AddInputField("Date", app.time.Format(dateLayout), 20, nil, nil).
		AddInputField("Time spent", "", 20, nil, nil).
		AddButton("Add", newWorklog).
		AddButton("Cancel", func() {
			app.ui.app.SetFocus(app.ui.dayView.latestIssuesList)
			app.ui.pages.RemovePage("worklog-form")
		})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.ui.pages.RemovePage("worklog-form")
			app.ui.app.SetFocus(app.ui.dayView.latestIssuesList)
			break
		}
		return event
	})
	form.SetBorder(true).SetTitle("New worklog").SetTitleAlign(tview.AlignLeft)
	_, _, pwidth, pheight := app.ui.grid.GetRect()
	formWidth := 36
	formHeight := 9
	form.SetRect(pwidth/2-(formWidth/2), pheight/2-3, formWidth, formHeight)
	app.ui.pages.AddPage("worklog-form", form, false, true)
	return form
}

func NewUpdateWorklogForm(d *DayView, workLogIssues []*WorkLogIssue, row int) *tview.Form {
	var form *tview.Form

	updateWorklog := func() {
		timeSpent := form.GetFormItem(0).(*tview.InputField).GetText()
		go func() {
			app.ui.loaderView.Show("Updating worklog...")
			defer app.ui.loaderView.Hide()
			err := workLogIssues[row].WorkLog.Update(timeSpent)
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
		go func() {
			app.ui.loaderView.Show("Deleting worklog...")
			defer app.ui.loaderView.Hide()
			err := app.workLogs.Delete(workLogIssues[row].WorkLog)
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
			app.ui.app.SetFocus(app.ui.dayView.worklogList)
		})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.ui.pages.RemovePage("worklog-form")
			app.ui.app.SetFocus(app.ui.dayView.worklogList)
			break
		}
		return event
	})
	form.SetBorder(true).SetTitle("Update worklog").SetTitleAlign(tview.AlignLeft)
	_, _, pwidth, pheight := app.ui.grid.GetRect()
	formWidth := 36
	formHeight := 7
	form.SetRect(pwidth/2-(formWidth/2), pheight/2-3, formWidth, formHeight)
	app.ui.pages.AddPage("worklog-form", form, false, true)
	return form
}
