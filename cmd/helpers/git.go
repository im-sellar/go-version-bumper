package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetCurrentGitBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func CommitChanges(currentBranchName string) error {
	cmd := exec.Command("git", "add", "package.json", "package-lock.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("error adding changes to staging area: %v", err)
	}

	if err != nil {
		return fmt.Errorf("error getting current git branch: %v", err)
	}

	if strings.Contains(currentBranchName, "feature") {
		cmd = exec.Command("git", "commit", "-m", "feat: upgrade version", "--no-verify")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else if strings.Contains(currentBranchName, "fix") {
		cmd = exec.Command("git", "commit", "-m", "fix: upgrade version", "--no-verify")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		return fmt.Errorf("branch name %s is not supported", currentBranchName)
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	return nil
}
