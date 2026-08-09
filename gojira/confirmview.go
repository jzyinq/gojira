package gojira

import (
	"github.com/rivo/tview"
)

type ConfirmView struct {
	*tview.Modal
	previousFocus tview.Primitive
	onConfirm     func()
}

func NewConfirmView() *ConfirmView {
	confirmView := &ConfirmView{tview.NewModal(), nil, nil}
	confirmView.SetTitle("Confirm")
	confirmView.SetText("Are you sure you want to proceed?")
	confirmView.AddButtons([]string{"Yes", "No"})
	confirmView.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" && confirmView.onConfirm != nil {
			confirmView.onConfirm()
		}
		app.ui.pages.HidePage("confirm")
		if confirmView.previousFocus != nil {
			app.ui.app.SetFocus(confirmView.previousFocus)
		}
	})
	app.ui.pages.AddPage("confirm", confirmView, true, false)
	return confirmView
}

func (e *ConfirmView) Confirm(confirmation string, previousFocus tview.Primitive, onConfirm func()) {
	e.SetText(confirmation)
	e.onConfirm = onConfirm
	e.previousFocus = previousFocus
	app.ui.pages.SendToFront("confirm")
	app.ui.pages.ShowPage("confirm")
}
