package gojira

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ErrorView struct {
	*tview.Modal
	previousFocus tview.Primitive
}

func NewErrorView() *ErrorView {
	errorView := &ErrorView{tview.NewModal(), nil}
	errorView.SetText("Something went wrong")
	errorView.SetTitle("Error!")
	errorView.SetBackgroundColor(tcell.ColorRed.TrueColor())
	errorView.AddButtons([]string{"OK"})
	errorView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			App.ui.pages.HidePage("error")
			if errorView.previousFocus != nil {
				App.ui.app.SetFocus(errorView.previousFocus)
			}
		}
		return event
	})
	App.ui.pages.AddPage("error", errorView, true, false)
	return errorView
}

func (e *ErrorView) ShowError(error string, previousFocus tview.Primitive) {
	e.previousFocus = previousFocus
	App.ui.pages.SendToFront("error")
	e.SetText(fmt.Sprintf("Error: %s", error))
	App.ui.pages.ShowPage("error")
	App.ui.app.SetFocus(e)
	App.ui.app.Draw()
}
