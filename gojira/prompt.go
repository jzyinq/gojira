package gojira

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"regexp"
	"strings"
)

func PromptForTimeSpent(promptLabel string) (string, error) {
	validate := func(input string) error {
		r, _ := regexp.Compile("^(([0-9]+)h)?\\s?(([0-9]+)m)?$")
		match := r.MatchString(input)
		if !match {
			return errors.New("Invalid timeSpent format - try 1h / 1h30m / 30m")
		}
		return nil
	}
	promptInput := promptui.Prompt{
		Label:    promptLabel,
		Validate: validate,
	}

	result, err := promptInput.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return FormatTimeSpent(TimeSpentToSeconds(result)), nil
}

func PromptForIssueSelection(issues []Issue) (Issue, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .Key }}?",
		Active:   "> {{ .Key | cyan }} {{ .Fields.Summary | green }}",
		Inactive: " {{ .Key | cyan }} {{ .Fields.Summary | blue }}",
		Selected: "{{ .Key | cyan }} {{ .Fields.Summary | blue }}",
		Details: `--------- Details ----------
{{ "Status:" | faint }}	{{ .Fields.Status.Name }}
`,
	}

	searcher := func(input string, index int) bool {
		issue := issues[index]
		name := strings.Replace(strings.ToLower(issue.Key+issue.Fields.Summary), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	promptSelect := promptui.Select{
		Label:     "Recently updated jira tickets assigned to you:",
		Items:     issues,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
		HideHelp:  true,
	}

	i, _, err := promptSelect.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return Issue{}, err
	}

	return issues[i], nil
}
