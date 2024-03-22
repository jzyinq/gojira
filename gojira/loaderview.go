package gojira

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
	"time"
)

type LoaderView struct {
	tview.Primitive
	text   *tview.TextView
	ctx    context.Context
	cancel context.CancelFunc
}

func NewLoaderView() *LoaderView {
	customModal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}
	text := tview.NewTextView()
	text.SetBorder(true)
	text.SetTextAlign(tview.AlignCenter)

	loaderView := &LoaderView{customModal(text, 36, 13), text, nil, nil}
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
					e.text.SetText(fmt.Sprintf("%s%s\n%s", GojiraAscii, msg, string(r)))
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
