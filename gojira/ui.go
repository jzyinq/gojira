package gojira

import (
	"github.com/rivo/tview"
)

type UserInteface struct {
	app       *tview.Application
	flex      *tview.Flex
	pages     *tview.Pages
	table     *tview.Table
	calendar  *Calendar
	summary   *Summary
	dayView   *DayView
	errorView *ErrorView
}

func newUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()

	// do I really need those declaration here
	app.ui.calendar = NewCalendar()
	app.ui.summary = NewSummary()
	app.ui.dayView = NewDayView()
	app.ui.errorView = NewErrorView()

	app.ui.flex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(app.ui.pages, 0, 5, true),
			0, 9, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(app.ui.summary.TextView, 0, 1, false).
			AddItem(app.ui.calendar.Table, 0, 10, false),
			0, 5, false)

	grid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0, 27).
		SetBorders(true)
	grid.SetTitle(" gojira ")
	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(app.ui.pages, 1, 0, 1, 2, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(app.ui.pages, 0, 0, 2, 1, 0, 100, false).
		AddItem(app.ui.summary, 0, 1, 1, 1, 0, 100, false).
		AddItem(app.ui.calendar, 1, 1, 1, 1, 0, 100, false)

	app.ui.flex.SetBorder(true).SetTitle(" gojira ")
	app.ui.app.SetRoot(grid, true).SetFocus(app.ui.pages) //FIXME set on proper item after rarrangements
}
