package gojira

import (
	"github.com/rivo/tview"
)

type UserInteface struct {
	app        *tview.Application
	flex       *tview.Flex
	pages      *tview.Pages
	table      *tview.Table
	calendar   *Calendar
	summary    *Summary
	dayView    *DayView
	errorView  *ErrorView
	LoaderView *LoaderView
}

func newUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()

	// do I really need those declaration here
	app.ui.calendar = NewCalendar()
	app.ui.summary = NewSummary()
	app.ui.dayView = NewDayView()
	app.ui.errorView = NewErrorView()
	app.ui.LoaderView = NewLoaderView()

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

	app.ui.app.SetRoot(grid, true).SetFocus(app.ui.pages) //FIXME set on proper item after rearrangements
}
