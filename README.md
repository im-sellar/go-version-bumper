# Go-Version-Bumper CLI App

This is a command-line application built using the Cobra library in Go. The application provides a `bump` command that upgrades the version of `package.json` and `package-lock.json` files and commits the changes to Git.

## Features

- Upgrade the version of `package.json` and `package-lock.json` files based on the current Git branch name.
- Commit the changes to Git with a commit message based on the branch name.

## Usage

To use the `bump` command, navigate to the root directory of your project and run:

```
version-bumper bump
```

This will upgrade the version of `package.json` and `package-lock.json` files based on the current Git branch name:

- If the branch name starts with "feature", the minor number will be upgraded.
- If the branch name starts with "fix", the patch number will be upgraded.

To commit the changes to Git, use the `--commit` or `-c` flag:

```
version-bumper bump --commit
```

This will add the modified `package.json` and `package-lock.json` files to the staging area and commit the changes to Git with a commit message based on the branch name:

- If the branch name starts with "feature", the commit message will be "feat: upgrade version".
- If the branch name starts with "fix", the commit message will be "fix: upgrade version".

## Installation

To install the application, clone the repository and build the application using Go:

```
git clone https://github.com/your-username/go-version-bumper.git
cd go-version-bumper
go build
```

This will create an executable file called `version-bumper` in the current directory. You can move this file to a directory in your `PATH` to make it accessible from anywhere.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
