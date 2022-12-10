package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

var (
	ctx          context.Context
	githubClient *github.Client
)

// init ctx and githubClient
func InitContextAndGitHubClient() error {
	ctx = context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.GitHubAPIToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)

	if params.EnterpriseURL == "" {
		githubClient = github.NewClient(httpClient)
	} else {
		var err error
		githubClient, err = github.NewEnterpriseClient(params.EnterpriseURL, params.EnterpriseURL, httpClient)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fetch the difference of commit IDs between develop and main
func FetchDiffCommitIDs() ([]string, error) {
	var diffCommitIDs []string
	pageNum := FirstPageNumberOfGitHubAPI

	for {
		perPage := PerPageDefault

		comparison, resp, err := githubClient.Repositories.CompareCommits(
			ctx,
			params.RepositoryOwner,
			params.RepositoryName,
			params.BaseBranchName,
			params.DevelopmentBranchName,
			&github.ListOptions{
				PerPage: perPage,
				Page:    pageNum,
			},
		)
		if err != nil {
			return diffCommitIDs, err
		}

		for _, commit := range comparison.Commits {
			diffCommitIDs = append(diffCommitIDs, commit.GetSHA())
		}

		if resp.NextPage == 0 {
			break
		} else {
			pageNum = resp.NextPage
		}
	}

	return diffCommitIDs, nil
}

// Fetch up to `limit` pull requests sorted by updated desc
func FetchClosedPullRequests(limit int) ([]*github.PullRequest, error) {
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
			return pullRequestsList, err
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

	return pullRequestsList, nil
}

// If the release pull request does not exist, create a new one, otherwise edit the title and body
func CreateOrUpdatePullRequest(title string, body string) (*github.PullRequest, bool, error) {
	var pullRequest *github.PullRequest
	isCreated := false

	releasePullRequests, _, err := githubClient.PullRequests.List(ctx, params.RepositoryOwner, params.RepositoryName, &github.PullRequestListOptions{
		Head:  fmt.Sprintf("%s:%s", params.RepositoryOwner, params.DevelopmentBranchName),
		Base:  params.BaseBranchName,
		State: "open",
	})
	if err != nil {
		return pullRequest, isCreated, err
	}

	if len(releasePullRequests) == 0 {
		pullRequest, err = createPullRequest(title, body)
		isCreated = true
	} else {
		pullRequest, err = updatePullRequest(title, body, releasePullRequests[0].GetNumber())
	}
	return pullRequest, isCreated, err
}

func createPullRequest(title string, body string) (*github.PullRequest, error) {
	newPullRequest, _, err := githubClient.PullRequests.Create(ctx, params.RepositoryOwner, params.RepositoryName, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Head:  &params.DevelopmentBranchName,
		Base:  &params.BaseBranchName,
	})

	return newPullRequest, err
}

func updatePullRequest(title string, body string, pullReqNumber int) (*github.PullRequest, error) {
	pullRequest, _, err := githubClient.PullRequests.Edit(ctx, params.RepositoryOwner, params.RepositoryName, pullReqNumber, &github.PullRequest{
		Title: &title,
		Body:  &body,
	})

	return pullRequest, err
}
