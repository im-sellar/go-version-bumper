/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	FeatureBranch = "feature"
	FixBranch     = "fix"
)

var ErrUnsupportedBranch = errors.New("unsupported branch name")

var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Upgrade the version of the package.json and package-lock.json",
	Long: `Upgrade the version of the package.json and package-lock.json and commit the changes to git.
The version is upgraded depending on the name of the branch. 
If the branch name begin with "feature" then the minor number will be upgraded.
If the branch name begin with "fix" then the patch number will be upgraded.`,

	Run: func(cmd *cobra.Command, args []string) {
		currentBranchName, err := getCurrentBranch()
		if err != nil {
			fmt.Printf("Error getting current git branch: %v\n", err)
			return
		}
		err = updateVersionInFile("./package.json", currentBranchName)
		if err != nil {
			fmt.Printf("Error updating package.json: %v\n", err)
			return
		}

		err = updateVersionInFile("./package-lock.json", currentBranchName)
		if err != nil {
			fmt.Printf("Error updating package-lock.json: %v\n", err)
			return
		}

		fmt.Println("Version updated successfully!")

		commit, err := cmd.Flags().GetBool("commit")
		if err != nil {
			fmt.Printf("Error getting commit flag: %v\n", err)
			return
		}

		if commit {
			err = commitChanges(currentBranchName)
			if err != nil {
				fmt.Printf("Error committing changes: %v\n", err)
				return
			}

			fmt.Println("Changes committed successfully!")
		}
	},
}

func updateVersionInFile(filePath string, currentBranchName string) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	currentVersion, err := extractVersion(fileContent)
	if err != nil {
		return err
	}

	newVersion, err := getNewVersionBasedOnBranch(currentVersion, currentBranchName)
	if err != nil {
		return err
	}

	updatedContent, err := updateVersionInContent(fileContent, newVersion)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func extractVersion(content []byte) (string, error) {
	re := regexp.MustCompile(`"version": "([^"]+)"`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) < 2 {
		return "", errors.New("error getting current version")
	}
	return matches[1], nil
}

func getNewVersionBasedOnBranch(currentVersion string, currentBranchName string) (string, error) {
	var index int
	if strings.Contains(currentBranchName, FeatureBranch) {
		index = 1
	} else if strings.Contains(currentBranchName, FixBranch) {
		index = 2
	} else {
		return "", ErrUnsupportedBranch
	}

	return upgradeVersion(currentVersion, index), nil
}

func updateVersionInContent(content []byte, newVersion string) (string, error) {
	re := regexp.MustCompile(`"version": "([^"]+)"`)
	replaced := false
	return re.ReplaceAllStringFunc(string(content), func(match string) string {
		if replaced {
			return match
		}
		replaced = true
		return `"version": "` + newVersion + `"`
	}), nil
}

func upgradeVersion(version string, index int) string {
	components := strings.Split(version, ".")

	if len(components) < 3 {
		return version
	}

	major, err := strconv.Atoi(components[0])
	if err != nil {
		return version
	}
	minor, err := strconv.Atoi(components[1])
	if err != nil {
		return version
	}
	patch, err := strconv.Atoi(components[2])
	if err != nil {
		return version
	}

	switch index {
	case 1:
		minor++
		patch = 0
	case 2:
		patch++
	default:
		return version
	}

	upgradedVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	return upgradedVersion
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func commitChanges(currentBranchName string) error {
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
func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.Flags().BoolP("commit", "c", false, "Commit the changes to git")
}
