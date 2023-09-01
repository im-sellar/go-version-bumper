package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func CalculateNewVersion(currentBranchName string) (string, error) {
	currentVersion, err := getCurrentVersion("./package.json")
	if err != nil {
		return "", err
	}

	var newVersion string
	if strings.Contains(currentBranchName, FeatureBranch) {
		newVersion = upgradeVersion(currentVersion, 1)
	} else if strings.Contains(currentBranchName, FixBranch) {
		newVersion = upgradeVersion(currentVersion, 2)
	} else {
		return "", fmt.Errorf("branch name %s is not supported", currentBranchName)
	}
	return newVersion, nil
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

func updateVersionFromFileContent(content []byte, newVersion string) string {
	re := regexp.MustCompile(`"version": "([^"]+)"`)
	replaced := false
	return re.ReplaceAllStringFunc(string(content), func(match string) string {
		if replaced {
			return match
		}
		replaced = true
		return `"version": "` + newVersion + `"`
	})
}
