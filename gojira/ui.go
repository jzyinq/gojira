package gojira

import (
	"github.com/rivo/tview"
)

type UserInteface struct {
	app        *tview.Application
	pages      *tview.Pages
	grid       *tview.Grid
	calendar   *Calendar
	summary    *Summary
	dayView    *DayView
	errorView  *ErrorView
	loaderView *LoaderView
}

func newUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()

	app.ui.calendar = NewCalendar()
	app.ui.summary = NewSummary()
	app.ui.dayView = NewDayView()
	app.ui.errorView = NewErrorView()
	app.ui.loaderView = NewLoaderView()

	app.ui.grid = tview.NewGrid().
		SetRows(1, 0, 0).
		SetColumns(0, 0, 27).
		SetBorders(true)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	app.ui.grid.AddItem(app.ui.pages, 0, 0, 2, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	app.ui.grid.AddItem(app.ui.pages, 0, 0, 3, 2, 0, 100, true).
		AddItem(app.ui.summary, 0, 2, 1, 1, 0, 100, false).
		AddItem(app.ui.calendar, 1, 2, 2, 1, 0, 100, false)
	app.ui.app.SetRoot(app.ui.grid, true).SetFocus(app.ui.pages)
}
