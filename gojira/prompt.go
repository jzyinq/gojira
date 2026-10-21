package gojira

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"regexp"
)

var timeSpentValidationRegex = regexp.MustCompile(`^(([0-9]+)h)?\s?(([0-9]+)m)?$`)

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

func newDescriptionInput(description *string) *huh.Input {
	return huh.NewInput().
		Title("Description").
		Placeholder("optional").
		Value(description)
}

func IssueWorklogForm(issues []Issue, worklogs []*Worklog) (Issue, string, string, error) {
	formOptions := make([]huh.Option[Issue], len(issues))
	for i, issue := range issues {
		timeSpent := ""
		worklog := findWorklogByIssueKey(worklogs, issue.Key)
		if worklog != nil {
			timeSpent = FormatTimeSpent(worklog.TimeSpentSeconds)
		}
		formOptions[i] = huh.NewOption(fmt.Sprintf("%-8s %-10s - %s", timeSpent, issue.Key, issue.Fields.Summary), issue)
	}
	chosenIssue := Issue{}
	timeSpent := ""
	description := ""
	descriptionInput := newDescriptionInput(&description)
	timeSpentInput := huh.NewInput().
		Title("Log time").
		Placeholder("1h / 1h30m / 30m").
		Value(&timeSpent).
		Validate(func(input string) error {
			match := timeSpentValidationRegex.MatchString(input)
			if !match {
				return errors.New("invalid timeSpent format - try 1h / 1h30m / 30m")
			}
			return nil
		})
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Issue]().
				Title("Choose issue").
				Description(fmt.Sprintf("Time logged for today: %s", FormatTimeSpent(CalculateTimeSpent(worklogs)))).
				Options(formOptions...).
				Value(&chosenIssue).
				Validate(func(issue Issue) error {
					// it's more like a "prepare next input" function
					timeSpent = ""
					description = ""
					worklog := findWorklogByIssueKey(worklogs, issue.Key)
					if worklog != nil {
						timeSpent = FormatTimeSpent(worklog.TimeSpentSeconds)
						description = worklog.Description
					}
					timeSpentInput.Value(&timeSpent)
					descriptionInput.Value(&description)
					timeSpentInput.Description(fmt.Sprintf("%s %s", chosenIssue.Key, chosenIssue.Fields.Summary))
					return nil
				}),
		),
		huh.NewGroup(timeSpentInput, descriptionInput),
	)
	form.WithTheme(huh.ThemeDracula())
	customizedKeyMap := huh.NewDefaultKeyMap()
	customizedKeyMap.Input.Prev = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back"))

	// merge NewDefaultKeyMap with custom keymap
	form.WithKeyMap(customizedKeyMap)
	err := form.Run()
	if err != nil {
		return chosenIssue, timeSpent, description, err
	}

	return chosenIssue, timeSpent, description, nil
}

func InputTimeSpentForm(issue Issue, timeSpent, description string) (string, string, error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Log time").
				Description(fmt.Sprintf("%s %s", issue.Key, issue.Fields.Summary)).
				Placeholder("1h / 1h30m / 30m").
				Value(&timeSpent).
				Validate(func(input string) error {
					match := timeSpentValidationRegex.MatchString(input)
					if !match {
						return errors.New("invalid timeSpent format - try 1h / 1h30m / 30m")
					}
					return nil
				}),
			newDescriptionInput(&description),
		),
	)
	form.WithTheme(huh.ThemeDracula())
	err := form.Run()
	if err != nil {
		return timeSpent, description, err
	}

	return timeSpent, description, nil
}
