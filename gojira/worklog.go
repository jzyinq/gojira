package gojira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

func Test() {
	test := tview.NewApplication()
	table := tview.NewTable().SetSelectable(true, false)
	frame := tview.NewFrame(table).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Worklogs", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("gojira v0.0.9", true, tview.AlignRight, tcell.ColorWhite).
		AddText("(p)revious day (n)ext day", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("(d)elete worklog (enter) update worklog", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("Status...", false, tview.AlignCenter, tcell.ColorGreen)

	headers := strings.Split("Key Description Worklog", " ")
	cols, rows := 3, 40
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			table.SetCell(r, c,
				tview.NewTableCell(headers[word]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(headers)
		}
	}

	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			test.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		// tutaj można wstawić input z wyedytowaniem workloga
		test.SetRoot(newWorklogForm(test, frame, table.GetCell(row, column).Text), true)
	})

	if err := test.SetRoot(frame, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	if err := test.SetRoot(frame, true).SetFocus(table).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func newWorklogForm(app *tview.Application, frame *tview.Frame, workLog string) *tview.Form {
	form := tview.NewForm().
		AddInputField("Worklog", workLog, 20, nil, nil).
		AddButton("Update", func() {
			app.SetRoot(frame, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(frame, true)
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	return form
}
