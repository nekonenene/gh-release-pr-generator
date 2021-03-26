# GitHub Release PullRequest Generator

This CLI app supports you to create a **release pull request**.


## Usage

### 1. Download the binary file

```sh
go install github.com/nekonenene/gh-release-pr-generator@latest
```

### 2. Get GitHub API Token

See [here](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token)

### 3. Run

```sh
gh-release-pr-generator --token <YOUR_GITHUB_TOKEN> --repo-owner <REPOSITORY_OWNER_NAME> --repo-name <REPOSITORY_NAME> --dev-branch <DEVELOPMENT_BRANCH_NAME> --prod-branch <PRODUCTION_BRANCH_NAME>
```

Example:

```sh
gh-release-pr-generator --token 123456789abcd123456789abcd123456789abcd --repo-owner nekonenene --repo-name my-repository-name --dev-branch staging --prod-branch production
```

You can see all parameters:

```sh
gh-release-pr-generator --help
```


## Build

```sh
go build -o bin/gh-release-pr-generator
```
