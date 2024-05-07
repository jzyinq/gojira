package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ErrorView struct {
	*tview.Modal
}

func NewErrorView() *ErrorView {
	errorView := &ErrorView{tview.NewModal()}
	errorView.SetText("Something went wrong")
	errorView.SetTitle("Error!")
	errorView.SetBackgroundColor(tcell.ColorRed.TrueColor())
	errorView.AddButtons([]string{"OK"})
	errorView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			app.ui.pages.HidePage("error")
		}
		return event
	})
	app.ui.pages.AddPage("error", errorView, true, false)
	return errorView
}

func (e *ErrorView) ShowError(error string) {
	app.ui.pages.SendToFront("error")
	e.SetText(fmt.Sprintf("Error: %s", error))
	app.ui.pages.ShowPage("error")
	app.ui.app.SetFocus(e)
	app.ui.app.Draw()
}
