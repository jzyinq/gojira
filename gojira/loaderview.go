package gojira

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"time"
)

type LoaderView struct {
	*tview.TextView
	ctx    context.Context
	cancel context.CancelFunc
}

func NewLoaderView() *LoaderView {
	text := tview.NewTextView()
	text.SetBorder(true)
	text.SetTextAlign(tview.AlignCenter)

	loaderView := &LoaderView{text, nil, nil}
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
	app.ui.pages.SendToFront("loader")
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				return
			default:
				for _, r := range `-\|/` {
					e.SetText(fmt.Sprintf("%s%s\n%s", GojiraAscii, msg, string(r)))
					app.ui.app.Draw()
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
	app.ui.pages.ShowPage("loader")
	app.ui.app.Draw()
}

func (e *LoaderView) Hide() {
	if e.cancel != nil {
		e.cancel()
	}
	app.ui.pages.HidePage("loader")
	app.ui.app.Draw()
}
