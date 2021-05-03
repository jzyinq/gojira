package gojira

import (
	"os/exec"
	"regexp"
)

func GetTicketFromGitBranch() string {
	gitBranch, err := exec.Command("git", "branch", "--show-current").CombinedOutput()
	if err != nil {
		return ""
	}
	branchName := string(gitBranch)
	r, _ := regexp.Compile("([A-Z]+-[0-9]+)")
	match := r.FindString(branchName)
	return match
}
