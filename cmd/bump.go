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

	branchName, err := getCurrentBranch()
	if err != nil {
		return fmt.Errorf("error getting current git branch: %v", err)
	}

	// Update the version based on the branch name
	var newVersion string
	if strings.Contains(branchName, "feature") {
		// Upgrade the minor number
		newVersion = upgradeVersion(currentVersion, 1)
	} else if strings.Contains(branchName, "fix") {
		// Upgrade the patch number
		newVersion = upgradeVersion(currentVersion, 2)
	} else {
		return fmt.Errorf("branch name %s is not supported", branchName)
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
func init() {
	rootCmd.AddCommand(bumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
