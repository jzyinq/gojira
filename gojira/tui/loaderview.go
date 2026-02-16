package gojira

import (
	"context"
	"fmt"
	"gojira/gojira"
	"time"

	"github.com/rivo/tview"
)

type LoaderView struct {
	*LoaderModal
	ctx    context.Context
	cancel context.CancelFunc
	text   string
}

func NewLoaderView() *LoaderView {
	loaderView := &LoaderView{NewModal(), nil, nil, ""}
	loaderView.SetBorder(false)
	loaderView.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	App.ui.pages.AddPage("loader", loaderView, true, false)
	return loaderView
}

// FIXME something is off here
func (e *LoaderView) Wrap(msg string, callable func()) {
	go func() {
		e.Show(msg)
		defer e.Hide()
		callable()
	}()
}

func (e *LoaderView) Show(msg string) {
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.UpdateText(msg)
	App.ui.pages.SendToFront("loader")
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				return
			default:
				for _, r := range `-\|/` {
					e.SetText(fmt.Sprintf("%s%s\n%s", gojira.AppAsciiArt, e.text, string(r)))
					App.ui.app.Draw()
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
	App.ui.pages.ShowPage("loader")
	App.ui.app.Draw()
}

func (e *LoaderView) UpdateText(msg string) {
	e.text = msg
}

func (e *LoaderView) Hide() {
	focusedPrimitive := App.ui.app.GetFocus()
	if e.cancel != nil {
		e.cancel()
	}
	App.ui.pages.HidePage("loader")
	if focusedPrimitive != e {
		App.ui.app.SetFocus(focusedPrimitive)
	}
	App.ui.app.Draw()
}
