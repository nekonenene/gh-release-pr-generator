# GitHub Release PullRequest Generator

This CLI app supports you to create a **release pull request**.

It fetches pull requests which merged into the development branch, and generates new pull request would merge into the production branch. If the pull request already exists, it updates the title and the body of that. This app will be convenient if your project follows git-flow.

## Usage

### 1. Install

Go 1.16+:

```sh
go install github.com/nekonenene/gh-release-pr-generator@latest
```

Otherwise:

```sh
go get github.com/nekonenene/gh-release-pr-generator@latest
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


## GitHub Actions

The following is an example of the config file for [GitHub Actions](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions).

```yml
name: Generate Release Pull Request
on:
  push:
    branches:
      - develop
jobs:
  gh-release-pr-generator:
    name: gh-release-pr-generator
    runs-on: ubuntu-20.04
    env:
      TZ: Asia/Tokyo
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: ^1.16.2
      - name: Install gh-release-pr-generator
        run: go install github.com/nekonenene/gh-release-pr-generator@latest
      - name: Run gh-release-pr-generator
        run: gh-release-pr-generator --token ${{ secrets.GITHUB_TOKEN }} --repo-owner ${{ github.repository_owner }} --repo-name ${{ github.event.repository.name }} --dev-branch develop --prod-branch main
```


## Build

```sh
make build
```
