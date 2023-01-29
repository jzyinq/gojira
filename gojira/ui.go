package gojira

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type UserInteface struct {
	app   *tview.Application
	frame *tview.Frame
	pages *tview.Pages
}

func NewUi() {
	app.ui.app = tview.NewApplication()
	app.ui.pages = tview.NewPages()
	app.ui.frame = tview.NewFrame(app.ui.pages)
	app.ui.app.SetRoot(app.ui.frame, true)
}

func NewWorkLogView(workLogs []WorkLogIssue) {
	//FIXMe wyświetlaj dane z worklogów
	table := tview.NewTable().SetSelectable(true, false)
	headers := strings.Split("Key Description Worklog", " ")
	cols, rows := 3, 10
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
			app.ui.app.Stop()
		}
	}).SetSelectedFunc(func(row, column int) {
		// akcja na edycje tabeli
	})
	app.ui.pages.AddPage("worklog-view", table, true, true)
	app.ui.frame.SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Worklogs", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("gojira v0.0.9", true, tview.AlignRight, tcell.ColorWhite).
		AddText("(p)revious day   (n)ext day", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("(d)elete worklog (enter) update worklog", true, tview.AlignLeft, tcell.ColorYellow).
		AddText("Status...", false, tview.AlignCenter, tcell.ColorGreen)
}
