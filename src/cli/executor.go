package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

var (
	ctx          context.Context
	githubClient *github.Client
)

// init ctx and githubClient
func initContextAndClient() {
	ctx = context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.GitHubAPIToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)
	githubClient = github.NewClient(httpClient)
}

func Exec() {
	ParseParameters()
	initContextAndClient()

	diffCommitIDs := fetchDiffCommitIDs()

	fmt.Println(diffCommitIDs)

	pulls := fetchPullRequests(FetchPullRequestsLimitDefault)

	fmt.Println(len(pulls))

	for _, pull := range pulls {
		fmt.Println(pull.GetTitle())
		fmt.Println(pull.GetUser().GetLogin())

		// Skip not merged pull requests
		if pull.GetMergedAt().IsZero() {
			fmt.Print("Not Merged!\n\n")
			continue
		}
		fmt.Println(pull.GetMergedAt())
		fmt.Println(pull.GetMergeCommitSHA())
		fmt.Println()
	}
}

// Fetch the difference of commit IDs between develop and main
func fetchDiffCommitIDs() []string {
	comparison, _, err := githubClient.Repositories.CompareCommits(
		ctx,
		params.RepositoryOwner,
		params.RepositoryName,
		params.ProductionBranchName,
		params.DevelopmentBranchName,
	)
	if err != nil {
		log.Fatal(err)
	}

	var diffCommitIDs []string
	for _, commit := range comparison.Commits {
		diffCommitIDs = append(diffCommitIDs, commit.GetSHA())
	}

	return diffCommitIDs
}

// Fetch up to `limit` pull requests sorted by updated desc
func fetchPullRequests(limit int) []*github.PullRequest {
	var pullRequestsList []*github.PullRequest
	pageNum := FirstPageNumberOfGitHubAPI

	for {
		perPage := PerPageDefault
		if limit < PerPageDefault {
			perPage = limit
		}

		listOptions := github.ListOptions{
			PerPage: perPage,
			Page:    pageNum,
		}

		pulls, resp, err := githubClient.PullRequests.List(ctx, params.RepositoryOwner, params.RepositoryName, &github.PullRequestListOptions{
			Base:        params.DevelopmentBranchName,
			State:       "closed",
			Sort:        "updated",
			Direction:   "desc",
			ListOptions: listOptions,
		})
		if err != nil {
			log.Fatal(err)
		}

		pullRequestsList = append(pullRequestsList, pulls...)
		limit = limit - len(pulls)
		if limit == 0 {
			break
		}

		if resp.NextPage == 0 {
			break
		} else {
			pageNum = resp.NextPage
		}
	}

	return pullRequestsList
}
