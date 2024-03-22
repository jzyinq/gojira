package gojira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Modal is a centered message window used to inform the user or prompt them
// for an immediate decision. It needs to have at least one button (added via
// [Modal.AddButtons]) or it will never disappear.
//
// See https://github.com/rivo/tview/wiki/Modal for an example.
type Modal struct {
	*tview.Box

	// The frame embedded in the modal.
	frame *tview.Frame

	// The form embedded in the modal's frame.
	textView *tview.TextView

	// The message text (original, not word-wrapped).
	text string

	// The text color.
	textColor tcell.Color

	// The optional callback for when the user clicked one of the buttons. It
	// receives the index of the clicked button and the button's label.
	done func(buttonIndex int, buttonLabel string)

	focus bool
}

// NewModal returns a new modal message window.
func NewModal() *Modal {
	m := &Modal{
		Box:       tview.NewBox().SetBorder(true),
		textColor: tview.Styles.PrimaryTextColor,
		textView:  tview.NewTextView(),
	}

	return m
}

// SetBackgroundColor sets the color of the modal frame background.
func (m *Modal) SetBackgroundColor(color tcell.Color) *Modal {
	return m
}

// SetTextColor sets the color of the message text.
func (m *Modal) SetTextColor(color tcell.Color) *Modal {
	m.textColor = color
	m.textView.SetTextColor(color)
	return m
}

// SetDoneFunc sets a handler which is called when one of the buttons was
// pressed. It receives the index of the button as well as its label text. The
// handler is also called when the user presses the Escape key. The index will
// then be negative and the label text an empty string.
func (m *Modal) SetDoneFunc(handler func(buttonIndex int, buttonLabel string)) *Modal {
	m.done = handler
	return m
}

// SetText sets the message text of the window. The text may contain line
// breaks but style tag states will not transfer to following lines. Note that
// words are wrapped, too, based on the final size of the window.
func (m *Modal) SetText(text string) *Modal {
	m.text = text
	m.textView.SetText(text)
	return m
}

// SetFocus shifts the focus to the button with the given index.
func (m *Modal) SetFocus(index int) *Modal {
	m.focus = true
	return m
}

// Focus is called when this primitive receives focus.
func (m *Modal) Focus(delegate func(p tview.Primitive)) {

}

// HasFocus returns whether or not this primitive has focus.
func (m *Modal) HasFocus() bool {
	return m.focus
}

// Draw draws this primitive onto the screen.
func (m *Modal) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3
	height := 12
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)
	// Draw the frame.
	m.Box.DrawForSubclass(screen, m)
	x, y, width, height = m.GetInnerRect()
	m.textView.SetRect(x, y, width, height)
	m.textView.SetTextAlign(tview.AlignCenter)
	m.textView.Draw(screen)
}
