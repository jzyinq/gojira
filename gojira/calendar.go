package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type Calendar struct {
	*tview.Table
	year  int
	month time.Month
	day   int
}

func NewCalendar() *Calendar {
	t := time.Now()

	calendar := &Calendar{
		Table: tview.NewTable(),
		year:  t.Year(),
		month: t.Month(),
		day:   t.Day(),
	}

	calendar.update()

	// Set up input handling
	calendar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			calendar.day--
		case tcell.KeyRight:
			calendar.day++
		case tcell.KeyUp:
			calendar.day -= 7
		case tcell.KeyDown:
			calendar.day += 7
		}

		// Handle day overflow/underflow
		if calendar.day < 1 {
			calendar.day = 1
		}
		if calendar.day > time.Date(calendar.year, calendar.month+1, 0, 0, 0, 0, 0, time.Local).Day() {
			calendar.day = time.Date(calendar.year, calendar.month+1, 0, 0, 0, 0, 0, time.Local).Day()
		}

		calendar.update()

		return event
	})

	return calendar
}

func (c *Calendar) update() {
	c.Clear()

	t := time.Date(c.year, c.month, 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(c.year, c.month+1, 0, 0, 0, 0, 0, time.Local).Day()

	// Weekdays
	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	for i, day := range weekdays {
		c.SetCell(0, i, tview.NewTableCell(day))
	}

	// Days
	week := 1
	for i := 1; i <= daysInMonth; i++ {
		dayOfWeek := int(t.Weekday()) - 1 // Weekday() returns 1 (Monday) to 7 (Sunday)
		if dayOfWeek < 0 {
			dayOfWeek = 6 // Sunday
		}

		cell := tview.NewTableCell(fmt.Sprintf("%d", i))
		if i == c.day {
			cell.SetTextColor(tcell.ColorRed)
		}
		c.SetCell(week, dayOfWeek, cell)

		if dayOfWeek == 6 {
			week++
		}

		t = t.AddDate(0, 0, 1)
	}
}

//
//func main() {
//	calendar := NewCalendar()
//
//	app := tview.NewApplication()
//	app.SetRoot(calendar, true)
//
//	if err := app.Run(); err != nil {
//		panic(err)
//	}
//}
