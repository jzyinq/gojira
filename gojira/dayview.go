package gojira

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const IssueKeyColumn = 0
const IssueSummaryColumn = 1
const TimeSpentColumn = 2

// validateTimeSpentInput is a tview-compatible validation function for time input fields
// It accepts input that matches the time format or is empty (for partial input during typing)
func validateTimeSpentInput(textToCheck string, lastChar rune) bool {
	// Allow empty input (user is still typing)
	if textToCheck == "" {
		return true
	}
	// Allow digits, 'h', 'm', and space characters for valid time format
	if lastChar == 'h' || lastChar == 'm' || lastChar == ' ' || (lastChar >= '0' && lastChar <= '9') {
		return true
	}
	return false
}

type DayView struct {
	worklogList        *tview.Table
	worklogStatus      *tview.TextView
	latestIssuesList   *tview.Table
	latestIssuesStatus *tview.TextView
	searchInput        *tview.InputField
}

func NewDayView() *DayView { //nolint:funlen
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
	dayView.searchInput = tview.NewInputField().SetLabel("(l)Latest | (/)Search: ").SetFieldWidth(60).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				go func() {
					dayView.SearchIssues(dayView.searchInput.GetText())
				}()
			}
			if key == tcell.KeyEscape {
				app.ui.app.SetFocus(dayView.latestIssuesList)
			}
		}).SetFieldStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack))

	dayView.worklogList.SetBorder(true)
	dayView.latestIssuesList.SetBorder(true)
	dayView.latestIssuesList.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGray))
	dayView.latestIssuesList.SetFocusFunc(func() {
		dayView.latestIssuesList.SetSelectedStyle(
			tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	})
	dayView.latestIssuesList.SetBlurFunc(func() {
		dayView.latestIssuesList.SetSelectedStyle(
			tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite))
	})
	dayView.worklogList.SetFocusFunc(func() {
		dayView.worklogList.SetSelectedStyle(
			tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	})
	dayView.worklogList.SetBlurFunc(func() {
		dayView.worklogList.SetSelectedStyle(
			tcell.StyleDefault.Background(tcell.ColorGrey).Foreground(tcell.ColorWhite))
	})
	dayView.worklogStatus.SetText(
		fmt.Sprintf("Worklogs - %s - [?h[white]]",
			app.time.Format("2006-01-02"),
		)).SetDynamicColors(true)

	flexView := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dayView.worklogStatus, 1, 1, false).
		AddItem(dayView.worklogList, 0, 1, true).
		AddItem(dayView.latestIssuesStatus, 1, 1, false).
		AddItem(dayView.latestIssuesList, 0, 1, false).
		AddItem(dayView.searchInput, 1, 1, false)

	dayView.worklogList.SetCell(0, IssueKeyColumn,
		tview.NewTableCell("Loading...").SetAlign(tview.AlignLeft),
	)

	// Make tab key able to switch between the two tables
	// Change focues table active row color to yellow and inactive to white
	flexView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if app.ui.app.GetFocus() == dayView.worklogList {
				app.ui.app.SetFocus(dayView.latestIssuesList)
				return nil
			}
			app.ui.app.SetFocus(dayView.worklogList)
			return nil
		}
		if event.Rune() == '/' {
			app.ui.app.SetFocus(dayView.searchInput)
			return nil
		}
		if event.Rune() == 'l' && app.ui.app.GetFocus() != dayView.searchInput {
			app.ui.loaderView.WithLoader("Searching...", func() error {
				dayView.loadLatest()
				return nil
			})
			return nil
		}
		return event
	})

	app.ui.pages.AddPage("worklog-view", flexView, true, true)

	dayView.worklogList.SetInputCapture(controlCalendar)
	dayView.latestIssuesList.SetInputCapture(controlCalendar)

	go func() {
		dayView.loadLatest()
	}()

	return dayView
}

var loadingWorklogs = make(chan bool, 1)

func loadWorklogs() {
	select {
	case loadingWorklogs <- true:
		go func() {
			defer func() { <-loadingWorklogs }()
			err := NewWorklogIssues()
			if err != nil {
				app.ui.errorView.ShowError(err.Error(), nil)
			}
			app.ui.dayView.update()
		}()
	default:
		// The goroutine is already loadingWorklogs, do nothing
	}
}

// updateAll updates the day view, calendar, and summary in one call
func (d *DayView) updateAll() {
	d.update()
	app.ui.pages.RemovePage("worklog-form")
	app.ui.calendar.update()
	app.ui.summary.update()
}

// populateIssuesList populates the latest issues list table with issues
func (d *DayView) populateIssuesList(issues []Issue) {
	d.latestIssuesList.Clear()
	d.latestIssuesList.SetSelectable(true, false)
	color := tcell.ColorWhite
	for r := 0; r < len(issues); r++ {
		d.latestIssuesList.SetCell(r, IssueKeyColumn,
			tview.NewTableCell(issues[r].Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.latestIssuesList.SetCell(r, IssueSummaryColumn,
			tview.NewTableCell(issues[r].Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	d.latestIssuesList.Select(0, IssueKeyColumn).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		NewAddWorklogForm(d, issues, row)
	})
}

func (d *DayView) update() {
	logs, _ := app.workLogsIssues.IssuesOnDate(app.time)
	d.worklogList.Clear()
	d.worklogList.SetSelectable(true, false)
	color := tcell.ColorWhite
	for r := 0; r < len(logs); r++ {
		d.worklogList.SetCell(r, IssueKeyColumn,
			tview.NewTableCell((logs)[r].Issue.Key).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.worklogList.SetCell(r, IssueSummaryColumn,
			tview.NewTableCell((logs)[r].Issue.Fields.Summary).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
		d.worklogList.SetCell(r, TimeSpentColumn,
			tview.NewTableCell(
				FormatTimeSpent((logs)[r].Worklog.TimeSpentSeconds)).SetTextColor(color).SetAlign(tview.AlignLeft),
		)
	}
	d.worklogList.Select(0, IssueKeyColumn).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		NewUpdateWorklogForm(d, logs, row)
	})
	d.worklogList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDelete:
			app.ui.loaderView.WithLoader("Deleting worklog...", func() error {
				row, _ := d.worklogList.GetSelection()
				if err := app.workLogs.Delete(logs[row].Worklog); err != nil {
					return err
				}
				d.updateAll()
				return nil
			})
		default:
		}
		return controlCalendar(event)
	})
	timeSpent := CalculateTimeSpent(getWorklogsFromWorklogIssues(logs))
	d.worklogStatus.SetText(
		fmt.Sprintf("Worklogs - %s - [%s%s[white]]",
			app.time.Format("2006-01-02"),
			GetTimeSpentColorTag(timeSpent, 8),
			FormatTimeSpent(timeSpent),
		)).SetDynamicColors(true)
}

func (d *DayView) loadLatest() {
	d.latestIssuesStatus.SetText("Latest issues").SetDynamicColors(true)
	issues, err := app.jiraClient.GetLatestIssues()
	if app.ui.errorView.ShowErrorIfPresent(err, nil) {
		return
	}
	d.populateIssuesList(issues.Issues)
}

func (d *DayView) SearchIssues(search string) {
	app.ui.loaderView.WithLoaderAndFocus("Searching...", d.searchInput, func() error {
		if search == "" {
			return nil
		}
		jql := fmt.Sprintf("text ~ \"%s\"", search)
		if FindIssueKeyInString(search) != "" {
			jql = fmt.Sprintf("(text ~ \"%s\" OR issuekey = \"%s\")", search, search)
		}
		issues, err := app.jiraClient.GetIssuesByJQL(
			fmt.Sprintf("%s ORDER BY updated DESC, created DESC", jql), 10,
		)
		if err != nil {
			return err
		}
		d.populateIssuesList(issues.Issues)
		d.latestIssuesStatus.SetText("Search results:")
		app.ui.app.SetFocus(d.latestIssuesList)
		return nil
	})
}

// DateRange is a struct for holding the start and end dates
type DateRange struct {
	StartDate    time.Time
	EndDate      time.Time
	NumberOfDays int
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

	// count number of dayhs between StartDate and EndDate
	numberOfDays := int(endDate.Sub(startDate).Hours() / 24)

	return DateRange{
		StartDate:    startDate,
		EndDate:      endDate,
		NumberOfDays: numberOfDays,
	}, nil
}

func NewAddWorklogForm(d *DayView, issues []Issue, row int) *tview.Form { //nolint:funlen
	var form *tview.Form

	newWorklog := func() {
		logTime := form.GetFormItem(0).(*tview.InputField).GetText()
		timeSpent := form.GetFormItem(1).(*tview.InputField).GetText()

		// Validate time format before submission
		if timeSpent == "" {
			app.ui.loaderView.WithLoader("Time spent cannot be empty", func() error {
				time.Sleep(2 * time.Second)
				return nil
			})
			return
		}
		if !timeSpentValidationRegex.MatchString(timeSpent) {
			app.ui.loaderView.WithLoader("Invalid timeSpent format - try 1h / 1h30m / 30m", func() error {
				time.Sleep(2 * time.Second)
				return nil
			})
			return
		}

		app.ui.loaderView.WithLoader("Adding worklog...", func() error {
			issue, err := app.jiraClient.GetIssue(issues[row].Key)
			if err != nil {
				return err
			}
			// TODO use ParseDateRange and LogWork for each day in range
			dateRange, err := ParseDateRange(logTime)
			if err != nil {
				return err
			}
			for day := dateRange.StartDate; day.Before(dateRange.EndDate.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
				err := issue.LogWork(&day, timeSpent)
				app.ui.loaderView.UpdateText(fmt.Sprintf("Adding worklog for %s ...", day.Format(dateLayout)))
				if err != nil {
					return err
				}
			}
			d.updateAll()
			return nil
		})
	}

	form = tview.NewForm().
		AddInputField("Date", app.time.Format(dateLayout), 20, nil, nil).
		AddInputField("Time spent", "", 20, validateTimeSpentInput, nil).
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
			return nil
		}
		return event
	})
	form.SetBorder(true).SetTitle("New worklog").SetTitleAlign(tview.AlignLeft)
	_, _, pwidth, pheight := app.ui.grid.GetRect()
	formWidth := 36
	formHeight := 9
	form.SetRect(pwidth/2-(formWidth/2), pheight/2-3, formWidth, formHeight)
	form.SetFocus(1)
	app.ui.pages.AddPage("worklog-form", form, false, true)
	return form
}

func NewUpdateWorklogForm(d *DayView, workLogIssues []*WorklogIssue, row int) *tview.Form { //nolint:funlen
	var form *tview.Form

	updateWorklog := func() {
		timeSpent := form.GetFormItem(0).(*tview.InputField).GetText()

		// Validate time format before submission (require non-empty for updates)
		if timeSpent == "" {
			app.ui.loaderView.WithLoader("Time spent cannot be empty", func() error {
				time.Sleep(2 * time.Second)
				return nil
			})
			return
		}
		if !timeSpentValidationRegex.MatchString(timeSpent) {
			app.ui.loaderView.WithLoader("Invalid timeSpent format - try 1h / 1h30m / 30m", func() error {
				time.Sleep(2 * time.Second)
				return nil
			})
			return
		}

		app.ui.loaderView.WithLoader("Updating worklog...", func() error {
			if err := workLogIssues[row].Worklog.Update(timeSpent); err != nil {
				return err
			}
			d.updateAll()
			return nil
		})
	}

	deleteWorklog := func() {
		app.ui.loaderView.WithLoader("Deleting worklog...", func() error {
			if err := app.workLogs.Delete(workLogIssues[row].Worklog); err != nil {
				return err
			}
			d.updateAll()
			return nil
		})
	}

	form = tview.NewForm().
		AddInputField("Time spent",
			FormatTimeSpent(workLogIssues[row].Worklog.TimeSpentSeconds), 20, validateTimeSpentInput, nil).
		AddButton("Update", updateWorklog).
		AddButton("Delete", deleteWorklog).
		AddButton("Cancel", func() {
			app.ui.pages.RemovePage("worklog-form")
			app.ui.app.SetFocus(app.ui.dayView.worklogList)
		})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDelete:
			deleteWorklog()
			return nil
		case tcell.KeyEscape:
			app.ui.pages.RemovePage("worklog-form")
			app.ui.app.SetFocus(app.ui.dayView.worklogList)
			return nil
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
