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

func NewUi(app interface{}) {
	// We'll need to pass the app instance and access its fields
	// For now, let's keep the global access but fix the import cycle
}
