package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

func main() {
	app := tview.NewApplication()

	// Create a grid to represent the days of the week.
	grid := tview.NewGrid().
		SetRows(1, 1, 1, 1, 1, 1, 1, 1).
		SetColumns(-1, -1, -1, -1, -1, -1, -1, -1).
		SetBorders(true)

	// Add headers for the days of the week.
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	for i, day := range days {
		grid.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(day), 0, i, 1, 1, 1, 0, false)
	}

	// Get the current day.
	currentDay := time.Now().Day()

	// Add days for the month of May 2023.
	// The month starts on a Monday and has 31 days.
	buttons := make([]*tview.Button, 31)
	for i := 1; i <= 31; i++ {
		row := (i)/7 + 1
		column := (i)%7 - 1
		if column < 0 {
			column = 6
			row--
		}

		// Create a button for each day.
		button := tview.NewButton(tview.Escape(fmt.Sprintf("%2d", i)))
		button.SetStyle(tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack))
		buttons[i-1] = button

		grid.AddItem(button, row, column, 1, 1, 0, 0, false)
	}
	// Set the grid as the root of the application.
	app.SetRoot(grid, true)
	// Handle key presses for navigation.
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Get the current focus.
		focus := app.GetFocus()

		// Find the button in the buttons slice.
		for i, button := range buttons {
			if focus == button {
				// Determine the new focus based on the key pressed.
				switch event.Key() {
				case tcell.KeyUp:
					if i >= 7 {
						app.SetFocus(buttons[i-7])
					}
				case tcell.KeyDown:
					if i+7 < len(buttons) {
						app.SetFocus(buttons[i+7])
					}
				case tcell.KeyLeft:
					if i > 0 {
						app.SetFocus(buttons[i-1])
					}
				case tcell.KeyRight:
					if i+1 < len(buttons) {
						app.SetFocus(buttons[i+1])
					}
				}
			}
		}

		return event
	})

	// Set the initial focus to the current day.
	// Subtract 1 because slice indices start at 0.
	if currentDay <= len(buttons) {
		app.SetFocus(buttons[currentDay-1])
	}

	// Run the application.
	if err := app.Run(); err != nil {
		panic(err)
	}
}
