/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type PackageJSON map[string]interface{}

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Upgrade the version of the package.json and package-lock.json",
	Long: `Upgrade the version of the package.json and package-lock.json and commit the changes to git.
The version is upgraded depending on the name of the branch. 
If the branch name begin with "feature" then the minor number will be upgraded.
If the branch name begin with "fix" then the patch number will be upgraded.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := updateVersionInFile("./package.json")
		if err != nil {
			fmt.Printf("Error updating package.json: %v\n", err)
			return
		}

		err = updateVersionInFile("./package-lock.json")
		if err != nil {
			fmt.Printf("Error updating package-lock.json: %v\n", err)
			return
		}

		fmt.Println("Version updated successfully!")

		// Check if the commit flag is set
		commit, err := cmd.Flags().GetBool("commit")
		if err != nil {
			fmt.Printf("Error getting commit flag: %v\n", err)
			return
		}

		// If the commit flag is set, commit the changes to git
		if commit {
			err = commitChanges()
			if err != nil {
				fmt.Printf("Error committing changes: %v\n", err)
				return
			}

			fmt.Println("Changes committed successfully!")
		}
	},
}

func updateVersionInFile(filePath string) error {
	// Read the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Get the current version
	re := regexp.MustCompile(`"version": "([^"]+)"`)
	matches := re.FindStringSubmatch(string(fileContent))
	if len(matches) < 2 {
		return fmt.Errorf("error getting current version")
	}
	currentVersion := matches[1]

	currentBranchName, err := getCurrentBranch()
	if err != nil {
		return fmt.Errorf("error getting current git branch: %v", err)
	}

	// Update the version based on the branch name
	var newVersion string
	if strings.Contains(currentBranchName, "feature") {
		// Upgrade the minor number
		newVersion = upgradeVersion(currentVersion, 1)
	} else if strings.Contains(currentBranchName, "fix") {
		// Upgrade the patch number
		newVersion = upgradeVersion(currentVersion, 2)
	} else {
		return fmt.Errorf("branch name %s is not supported", currentBranchName)
	}

	// Update the file content with the new version
	replaced := false
	updatedFileContent := re.ReplaceAllStringFunc(string(fileContent), func(match string) string {
		if replaced {
			return match
		}
		replaced = true
		return `"version": "` + newVersion + `"`
	})

	// Write the updated content to the file
	err = os.WriteFile(filePath, []byte(updatedFileContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func upgradeVersion(version string, index int) string {
	// Split the version string into major, minor, and patch components
	components := strings.Split(version, ".")

	// Ensure that the version string has at least three components
	if len(components) < 3 {
		return version
	}

	// Parse the major, minor, and patch components
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

	// Upgrade the version based on the specified index
	switch index {
	case 1: // Upgrade minor number
		minor++
		patch = 0
	case 2: // Upgrade patch number
		patch++
	default:
		return version
	}

	// Construct the upgraded version string
	upgradedVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	return upgradedVersion
}

// getCurrentBranch returns the name of the current Git branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// commitChanges commits the changes to git
func commitChanges() error {
	// Add the changes to the staging area
	cmd := exec.Command("git", "add", "package.json", "package-lock.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding changes to staging area: %v", err)
	}

	currentBranchName, err := getCurrentBranch()
	if err != nil {
		return fmt.Errorf("error getting current git branch: %v", err)
	}

	// Commit the changes based on the branch name
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
