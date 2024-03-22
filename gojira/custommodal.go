package gojira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// LoaderModal is a centered message window used to inform
// the user about asynchronous action taking place.
//
// See https://github.com/rivo/tview/wiki/Modal for an example.
type LoaderModal struct {
	*tview.Box
	// The form embedded in the modal's frame.
	textView *tview.TextView
	// The message text (original, not word-wrapped).
	focus bool
}

// NewModal returns a new modal message window.
func NewModal() *LoaderModal {
	m := &LoaderModal{
		Box:      tview.NewBox().SetBorder(true),
		textView: tview.NewTextView(),
		focus:    false,
	}
	m.textView.SetTextColor(tview.Styles.PrimaryTextColor)
	m.textView.SetBorder(true)
	m.textView.SetTextAlign(tview.AlignCenter)
	return m
}

// SetBackgroundColor sets the color of the modal frame background.
func (m *LoaderModal) SetBackgroundColor(color tcell.Color) *LoaderModal {
	return m
}

// SetTextColor sets the color of the message text.
func (m *LoaderModal) SetTextColor(color tcell.Color) *LoaderModal {
	m.textView.SetTextColor(color)
	return m
}

// SetText sets the message text of the window. The text may contain line
// breaks but style tag states will not transfer to following lines. Note that
// words are wrapped, too, based on the final size of the window.
func (m *LoaderModal) SetText(text string) *LoaderModal {
	m.textView.SetText(text)
	return m
}

// SetFocus shifts the focus to the button with the given index.
func (m *LoaderModal) SetFocus(index int) *LoaderModal {
	m.focus = true
	return m
}

// Focus is called when this primitive receives focus.
func (m *LoaderModal) Focus(delegate func(p tview.Primitive)) {
	m.focus = true
}

// HasFocus returns whether or not this primitive has focus.
func (m *LoaderModal) HasFocus() bool {
	return m.focus
}

func (m *LoaderModal) Blur() {
	m.focus = false
}

// Draw draws this primitive onto the screen.
func (m *LoaderModal) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()
	// FIXME should be configurable or dynamic
	width := 38
	height := 14
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)
	m.Box.DrawForSubclass(screen, m)
	x, y, width, height = m.GetInnerRect()
	m.textView.SetRect(x, y, width, height)
	m.textView.Draw(screen)
}
