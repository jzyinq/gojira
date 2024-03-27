package gojira

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"time"
)

type Calendar struct {
	*tview.Table
	text  *tview.TextView
	year  int
	month time.Month
	day   int
}

func NewCalendar() *Calendar {
	t := app.time

	calendar := &Calendar{
		Table: tview.NewTable(),
		text:  tview.NewTextView().SetText("Calendar"),
		year:  t.Year(),
		month: t.Month(),
		day:   t.Day(),
	}

	calendar.update()

	// calendar browsing is pinned to dayview

	return calendar
}

func (c *Calendar) update() {
	c.day = app.time.Day()
	c.month = app.time.Month()
	c.year = app.time.Year()
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

		cell := tview.NewTableCell(fmt.Sprintf("%d", i)).SetAlign(tview.AlignCenter)

		calendarDay := time.Date(c.year, c.month, i, 0, 0, 0, 0, time.UTC)
		if calendarDay.Before(time.Now().Local()) {
			cell.SetBackgroundColor(tcell.ColorGray)
		}
		if len(app.workLogs.logs) > 0 {
			worklogs, err := app.workLogs.LogsOnDate(&calendarDay)
			if err != nil {
				panic(err)
			}
			timeSpent := CalculateTimeSpent(worklogs)
			color := GetTimeSpentColor(timeSpent, 8)
			cell.SetTextColor(color)
			if (dayOfWeek == 5 || dayOfWeek == 6) && timeSpent == 0 {
				cell.SetTextColor(tcell.ColorGrey)
				if calendarDay.Before(time.Now().Local()) {
					cell.SetTextColor(tcell.ColorBlack)
				}
			}
		}
		if i == c.day {
			cell.SetTextColor(tcell.ColorWhite)
			cell.SetBackgroundColor(tcell.ColorDimGray)
		}
		c.SetCell(week, dayOfWeek, cell)

		if dayOfWeek == 6 {
			week++
		}

		t = t.AddDate(0, 0, 1)
	}
}

func controlCalendar(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyLeft, tcell.KeyRight:
		timePeriod := -time.Hour * 24
		if event.Key() == tcell.KeyRight {
			timePeriod = time.Hour * 24
		}
		newTime := app.time.Add(timePeriod)
		logrus.Debug("Changing date to ", newTime)
		app.time = &newTime
		loadWorklogs()
		app.ui.calendar.update()
	}
	return event
}
