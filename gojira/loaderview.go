package gojira

import (
	"context"
	"fmt"
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
	app.ui.pages.AddPage("loader", loaderView, true, false)
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

// WithLoader wraps a function with loader show/hide, handling errors automatically
// If the function returns an error, it displays it via the error view
func (e *LoaderView) WithLoader(msg string, fn func() error) {
	go func() {
		e.Show(msg)
		defer e.Hide()
		if err := fn(); err != nil {
			app.ui.errorView.ShowError(err.Error(), nil)
		}
	}()
}

// WithLoaderAndFocus wraps a function with loader show/hide, handling errors with custom focus
func (e *LoaderView) WithLoaderAndFocus(msg string, focus tview.Primitive, fn func() error) {
	go func() {
		e.Show(msg)
		defer e.Hide()
		if err := fn(); err != nil {
			app.ui.errorView.ShowError(err.Error(), focus)
		}
	}()
}

func (e *LoaderView) Show(msg string) {
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.UpdateText(msg)
	app.ui.pages.SendToFront("loader")
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				return
			default:
				for _, r := range `-\|/` {
					e.SetText(fmt.Sprintf("%s%s\n%s", AppAsciiArt, e.text, string(r)))
					app.ui.app.Draw()
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
	app.ui.pages.ShowPage("loader")
	app.ui.app.Draw()
}

func (e *LoaderView) UpdateText(msg string) {
	e.text = msg
}

func (e *LoaderView) Hide() {
	focusedPrimitive := app.ui.app.GetFocus()
	if e.cancel != nil {
		e.cancel()
	}
	app.ui.pages.HidePage("loader")
	if focusedPrimitive != e {
		app.ui.app.SetFocus(focusedPrimitive)
	}
	app.ui.app.Draw()
}
