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
	loaderView *LoaderView
}

func newUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()

	// do I really need those declaration here
	app.ui.calendar = NewCalendar()
	app.ui.summary = NewSummary()
	app.ui.dayView = NewDayView()
	app.ui.errorView = NewErrorView()
	app.ui.loaderView = NewLoaderView()

	//customModal := func(p tview.Primitive, width, height int) tview.Primitive {
	//	return tview.NewGrid().
	//		SetColumns(0, width, 0).
	//		SetRows(0, height, 0).
	//		AddItem(p, 1, 1, 1, 1, 0, 0, true)
	//}
	grid := tview.NewGrid().
		SetRows(1, 0, 0).
		SetColumns(0, 0, 27).
		SetBorders(true)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(app.ui.pages, 0, 0, 2, 3, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(app.ui.pages, 0, 0, 3, 2, 0, 100, false).
		AddItem(app.ui.summary, 0, 2, 1, 1, 0, 100, false).
		AddItem(app.ui.calendar, 1, 2, 2, 1, 0, 100, false)
	app.ui.app.SetRoot(grid, true).SetFocus(app.ui.pages) //FIXME set on proper item after rearrangements
}
