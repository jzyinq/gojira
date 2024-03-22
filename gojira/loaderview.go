package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LoaderView struct {
	*tview.Modal
}

func NewLoaderView() *LoaderView {
	errorView := &LoaderView{tview.NewModal()}
	errorView.SetText("Loader place...holder! xddd")
	errorView.SetTitle("Doing something...")
	errorView.SetBackgroundColor(tcell.ColorBlue.TrueColor())
	app.ui.pages.AddPage("loader", errorView, true, false)
	return errorView
}

func (e *LoaderView) Wrap(msg string, callable func()) {
	go func() {
		e.Show(msg)
		defer e.Hide()
		callable()
	}()
}

func (e *LoaderView) Show(msg string) {
	app.ui.pages.SendToFront("loader")
	e.SetText(fmt.Sprintf("Hi! I'm the loader - please wait...\n%s", msg))
	app.ui.pages.ShowPage("loader")
	app.ui.app.Draw()
}

func (e *LoaderView) Hide() {
	app.ui.pages.HidePage("loader")
	app.ui.app.Draw()
}
