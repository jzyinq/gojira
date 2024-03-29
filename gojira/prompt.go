package gojira

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/sirupsen/logrus"
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

func IssueWorklogForm(issues []Issue) (Issue, string, error) {
	formOptions := make([]huh.Option[Issue], len(issues))
	for i, issue := range issues {
		formOptions[i] = huh.NewOption(fmt.Sprintf("%s - %s", issue.Key, issue.Fields.Summary), issue)
	}
	chosenIssue := Issue{}
	timeSpent := ""
	timeSpentInput := huh.NewInput().
		Title("Log time").
		Placeholder("1h / 1h30m / 30m").
		Value(&timeSpent).
		Validate(func(input string) error {
			r, _ := regexp.Compile(`^(([0-9]+)h)?\s?(([0-9]+)m)?$`)
			match := r.MatchString(input)
			if !match {
				return errors.New("invalid timeSpent format - try 1h / 1h30m / 30m")
			}
			return nil
		})
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Issue]().
				Title("Choose issue").
				Description(fmt.Sprintf("Time logged for today: %s", FormatTimeSpent(CalculateTimeSpent(app.workLogs.logs)))).
				Options(formOptions...).
				Value(&chosenIssue).Validate(func(issue Issue) error {
				logrus.Info("Validating issue...")
				timeSpent = ""
				worklog := findWorklogByIssueKey(issue.Key)
				if worklog != nil {
					timeSpent = FormatTimeSpent(worklog.TimeSpentSeconds)
					logrus.Info("Found worklog, setting initial timeSpent to ", timeSpent)
				}
				timeSpentInput.Value(&timeSpent)
				timeSpentInput.Description(fmt.Sprintf("%s %s", chosenIssue.Key, chosenIssue.Fields.Summary))
				return nil
			}),
		),
		huh.NewGroup(timeSpentInput),
	)
	form.WithTheme(huh.ThemeDracula())
	err := form.Run()
	if err != nil {
		return chosenIssue, timeSpent, err
	}

	return chosenIssue, timeSpent, nil
}

func InputTimeSpentForm(issue Issue, timeSpent string) (string, error) {
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
