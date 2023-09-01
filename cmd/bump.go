package cmd

import (
	"fmt"
	"github.com/im-sellar/go-version-bumper/cmd/helpers"
	"github.com/spf13/cobra"
)

var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Upgrade the version of the package.json and package-lock.json",
	Long: `Upgrade the version of the package.json and package-lock.json and commit the changes to git.
The version is upgraded depending on the name of the branch. 
If the branch name begin with "feature" then the minor number will be upgraded.
If the branch name begin with "fix" then the patch number will be upgraded.`,

	Run: func(cmd *cobra.Command, args []string) {
		currentBranchName, err := helpers.GetCurrentGitBranch()
		if err != nil {
			fmt.Printf("Error getting current git branch: %v\n", err)
			return
		}

		newVersion, err := helpers.CalculateNewVersion(currentBranchName)
		if err != nil {
			fmt.Printf("Error setting new version: %v\n", err)
			return
		}

		err = helpers.UpdatePackageFile("./package.json", newVersion)
		if err != nil {
			fmt.Printf("Error updating package.json: %v\n", err)
			return
		}

		err = helpers.UpdatePackageFile("./package-lock.json", newVersion)
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
			err = helpers.CommitChanges(currentBranchName)
			if err != nil {
				fmt.Printf("Error committing changes: %v\n", err)
				return
			}

			fmt.Println("Changes committed successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.Flags().BoolP("commit", "c", false, "Commit the changes to git")
}
