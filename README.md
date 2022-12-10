# GitHub Release PullRequest Generator

This CLI app supports you to create a **release pull request**.

It fetches pull requests which merged into the development branch, and generates new pull request would merge into the production branch. If the pull request already exists, it updates the title and the body of that. This app will be convenient if your project follows git-flow.

<p align="center">
  <img width="80%" alt="Screenshot of the Release Pull Request" src="https://user-images.githubusercontent.com/11713748/112718623-90307480-8f37-11eb-8139-a9bbf9b81ab1.png">
</p>


## Installation

Go 1.16+:

```sh
go install github.com/nekonenene/gh-release-pr-generator@latest
```

Otherwise:

```sh
go get github.com/nekonenene/gh-release-pr-generator@latest
```


## Usage

First, you need to get GitHub API Token to control your repository, please see [here](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token).

### Example

```sh
gh-release-pr-generator --token 123456789abcd123456789abcd --repo-owner nekonenene --repo-name my-repository-name --head-branch staging --base-branch production
```

### Parameters

You can see all parameters:

```sh
gh-release-pr-generator --help
```

| Parameter | Description | Required? |
|:---:|:---:|:---:|
|-token| GitHub API Token | YES |
|-repo-owner| Repository owner name | YES |
|-repo-name| Repository name | YES |
|-base-branch<br>(-prod-branch)| Production branch name (default: `main`) |  |
|-head-branch<br>(-dev-branch)| Development branch name (default: `develop`) |  |
|-template-path| PATH of the [template file](#template-file) |  |
|-limit| Limit number of fetching pull requests (default: `100`) |  |
|-enterprise-url| URL of GitHub Enterprise (ex. https://github.your.domain ) |  |

### Template File

You can customize the title and the body of the release pull request. Create a template file and specific that with `-template-path` option.

The first line of the template will be the title, and the rest will be the body.

#### Example

```
Release {{ .Year }}-{{ .Month }}-{{ .Date }} {{ .Hour }}:{{ .Minute }}
# Pull Requests
{{ range $i, $pull := .Pulls }}
* {{ $pull.Title }} (#{{ $pull.Number }}) @{{ $pull.User.Login }}
{{- end }}
```

#### Parameters for Template File

| Paramter | Description |
|:---:|:---:|
| Year | Current Year (ex. `2006`) |
| YearShort | Current Year (ex. `06`) |
| Month | Current Month (ex. `01`) |
| MonthShort | Current Month (ex. `1`) |
| Date | Current Date (ex. `02`) |
| DateShort | Current Date (ex. `2`) |
| Weekday | Current Weekday (ex. `Monday`) |
| WeekdayShort | Current Weekday (ex. `Mon`) |
| Hour | Current Hour (ex. `03`) |
| HourShort | Current Hour (ex. `3`) |
| Minute | Current Minute (ex. `04`) |
| MinuteShort | Current Minute (ex. `4`) |
| Second | Current Second (ex. `05`) |
| SecondShort | Current Second (ex. `5`) |
| Time | Current [Time](https://golang.org/pkg/time/#Time) |
| Pulls | [Pull Requests](https://github.com/google/go-github/blob/master/github/pulls.go) Array |


### Run with GitHub Actions

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
    runs-on: ubuntu-latest
    env:
      TZ: Asia/Tokyo
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: ^1.19.4
      - name: Install gh-release-pr-generator
        run: go install github.com/nekonenene/gh-release-pr-generator@latest
      - name: Run gh-release-pr-generator
        run: gh-release-pr-generator --token ${{ secrets.GITHUB_TOKEN }} --repo-owner ${{ github.repository_owner }} --repo-name ${{ github.event.repository.name }} --dev-branch develop --prod-branch main
```


## Build

```sh
make build
```


## License

[MIT](https://choosealicense.com/licenses/mit/)
