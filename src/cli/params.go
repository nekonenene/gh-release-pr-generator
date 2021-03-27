package cli

import (
	"flag"
	"log"
)

type parameters struct {
	GitHubAPIToken        string
	RepositoryOwner       string
	RepositoryName        string
	ProductionBranchName  string
	DevelopmentBranchName string
	TemplatePath          string
}

var params parameters

func ParseParameters() {
	flag.StringVar(&params.GitHubAPIToken, "token", "", "[Required] GitHub API Token")
	flag.StringVar(&params.RepositoryOwner, "repo-owner", "", "[Required] Repository owner")
	flag.StringVar(&params.RepositoryName, "repo-name", "", "[Required] Repository name")
	flag.StringVar(&params.ProductionBranchName, "prod-branch", ProductionBranchNameDefault, "[Opiton] production branch name")
	flag.StringVar(&params.DevelopmentBranchName, "dev-branch", DevelopmentBranchNameDefault, "[Opiton] development branch name")
	flag.StringVar(&params.TemplatePath, "template-path", "", "[Opiton] template path for customizing the title and the body of the release pull request")
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
