package gojira

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"regexp"
	"strings"
	"text/template"
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
		Size:      5,
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

func PromptForWorkLogSelection(workLogIssues []WorkLogIssue) (*WorkLog, error) {
	//add timeSpent to available template functions
	funcMap := template.FuncMap{"timeSpent": FormatTimeSpent}
	//preserve previous functions
	for k, v := range promptui.FuncMap {
		funcMap[k] = v
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ .Issue.Key }}",
		Active:   "-> {{ .Issue.Key | cyan }} [{{ .WorkLog.TimeSpentSeconds | timeSpent }}] {{ .Issue.Fields.Summary }}",
		Inactive: "  {{ .Issue.Key | cyan }} [{{ .WorkLog.TimeSpentSeconds | timeSpent }}] {{ .Issue.Fields.Summary }}",
		Selected: "Selected work log {{ .Issue.Key | cyan }} [{{ .WorkLog.TimeSpentSeconds | timeSpent }}]",
		FuncMap:  funcMap,
	}

	searcher := func(input string, index int) bool {
		issue := workLogIssues[index]
		name := strings.Replace(strings.ToLower(issue.Issue.Key+issue.Issue.Fields.Summary), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	promptSelect := promptui.Select{
		Label:     "Today's work logs [" + CalculateTimeSpent(getWorkLogsFromWorkLogIssues(workLogIssues)) + "]",
		Items:     workLogIssues,
		Templates: templates,
		Size:      5,
		Searcher:  searcher,
		HideHelp:  true,
	}

	i, _, err := promptSelect.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return &WorkLog{}, err
	}

	return &workLogIssues[i].WorkLog, nil
}
