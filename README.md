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
gh-release-pr-generator --token 123456789abcd123456789abcd --repo-owner nekonenene --repo-name my-repository-name --dev-branch staging --prod-branch production
```

### Parameters

You can see all parameters:

```sh
gh-release-pr-generator --help
```

| Parameter | Description | Required? |
|:---:|:---:|:---:|
|-token| GitHub API Token | YES |
|-repo-owner| Repository Owner Name | YES |
|-repo-name| Repository Name | YES |
|-prod-branch| Production Branch Name (default: `main`) |  |
|-dev-branch| Development Branch Name (default: `develop`) |  |
|-template-path| PATH of the [Template File](#template-file) |  |

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


## License

[MIT](https://choosealicense.com/licenses/mit/)
