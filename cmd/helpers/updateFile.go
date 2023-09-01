package helpers

import (
	"fmt"
	"os"
	"regexp"
)

func UpdatePackageFile(filePath string, newVersion string) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	updatedContent := updateVersionFromFileContent(fileContent, newVersion)

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func getCurrentVersion(filePath string) (string, error) {
	// Read the file
	fileContent, err := os.ReadFile(filePath)

	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	// Get the current version
	re := regexp.MustCompile(`"version": "([^"]+)"`)
	matches := re.FindStringSubmatch(string(fileContent))
	if len(matches) < 2 {
		return "", fmt.Errorf("error getting current version")
	}
	return matches[1], nil
}
