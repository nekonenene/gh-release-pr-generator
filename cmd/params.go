package cmd

import (
	"flag"
	"log"
)

type parameters struct {
	GitHubAPIToken         string
	RepositoryOwner        string
	RepositoryName         string
	BaseBranchName         string
	DevelopmentBranchName  string
	TemplatePath           string
	FetchPullRequestsLimit int
	EnterpriseURL          string
	DryRun                 bool
}

var params parameters

func ParseParameters() {
	flag.StringVar(&params.GitHubAPIToken, "token", "", "[Required] GitHub API Token")
	flag.StringVar(&params.RepositoryOwner, "repo-owner", "", "[Required] Repository owner")
	flag.StringVar(&params.RepositoryName, "repo-name", "", "[Required] Repository name")
	flag.StringVar(&params.BaseBranchName, "base-branch", BaseBranchNameDefault, "[Opiton] base branch name")
	flag.StringVar(&params.BaseBranchName, "prod-branch", BaseBranchNameDefault, "[Opiton] alias of \"base-branch\"")
	flag.StringVar(&params.DevelopmentBranchName, "head-branch", DevelopmentBranchNameDefault, "[Opiton] head branch (development branch) name")
	flag.StringVar(&params.DevelopmentBranchName, "dev-branch", DevelopmentBranchNameDefault, "[Opiton] alias of \"head-branch\"")
	flag.StringVar(&params.TemplatePath, "template-path", "", "[Opiton] template path for customizing the title and the body of the release pull request")
	flag.IntVar(&params.FetchPullRequestsLimit, "limit", FetchPullRequestsLimitDefault, "[Opiton] limit number of fetching pull requests")
	flag.StringVar(&params.EnterpriseURL, "enterprise-url", "", "[Opiton] URL of GitHub Enterprise (ex. https://github.your.domain )")
	flag.BoolVar(&params.DryRun, "dry-run", false, "[Opiton] Only display, not create or update a pull request")
	flag.Parse()

	if params.GitHubAPIToken == "" {
		log.Fatalln("-token is required")
	}
	if params.RepositoryOwner == "" {
		log.Fatalln("-repo-owner is required")
	}
	if params.RepositoryName == "" {
		log.Fatalln("-repo-name is required")
	}
}
