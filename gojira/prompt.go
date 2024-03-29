package gojira

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"regexp"
)

func SelectActionForm(actions []string) (string, error) {
	formOptions := make([]huh.Option[string], len(actions))
	for i, action := range actions {
		formOptions[i] = huh.NewOption(action, action)
	}
	chosenAction := ""

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose action").
				Options(formOptions...).
				Value(&chosenAction),
		),
	)
	form.WithTheme(huh.ThemeDracula())
	err := form.Run()
	if err != nil {
		return chosenAction, err
	}

	return chosenAction, nil
}

func SelectIssueForm(issues []Issue) (Issue, error) {
	formOptions := make([]huh.Option[Issue], len(issues))
	for i, issue := range issues {
		formOptions[i] = huh.NewOption(fmt.Sprintf("%s - %s", issue.Key, issue.Fields.Summary), issue)
	}
	chosenIssue := Issue{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Issue]().
				Title("Choose issue").
				Description(fmt.Sprintf("Time logged for today: %s", FormatTimeSpent(CalculateTimeSpent(app.workLogs.logs)))).
				Options(formOptions...).
				Value(&chosenIssue),
		),
	)
	form.WithTheme(huh.ThemeDracula())
	err := form.Run()
	if err != nil {
		return chosenIssue, err
	}

	return chosenIssue, nil
}

func InputTimeSpentForm(issue Issue) (string, error) {
	timeSpent := ""

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Log time").
				Description(fmt.Sprintf("%s %s", issue.Key, issue.Fields.Summary)).
				Placeholder("1h / 1h30m / 30m").
				Value(&timeSpent).
				Validate(func(input string) error {
					r, _ := regexp.Compile(`^(([0-9]+)h)?\s?(([0-9]+)m)?$`)
					match := r.MatchString(input)
					if !match {
						return errors.New("Invalid timeSpent format - try 1h / 1h30m / 30m")
					}
					return nil
				}),
		),
	)
	form.WithTheme(huh.ThemeDracula())
	err := form.Run()
	if err != nil {
		return timeSpent, err
	}

	return timeSpent, nil
}
