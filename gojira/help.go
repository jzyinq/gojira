package gojira

import "github.com/rivo/tview"

type Help struct {
	*tview.TextView
}

func NewHelpView() *Help {
	help := &Help{
		TextView: tview.NewTextView(),
	}
	help.SetText("Shortcuts:\n\n" +
		"<- - previous day\n" +
		"-> - next day\n" +
		"shift + -> - next month\n" +
		"shift + <- - previous month\n" +
		"/ - search\n" +
		"del - delete log\n" +
		"enter - add/edit log\n" +
		"tab - switch between lists",
	)
	return help
}
