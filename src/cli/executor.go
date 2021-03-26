package cli

import (
	"context"
	"fmt"
	"log"
	"time"

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

	diffCommitIDs, err := fetchDiffCommitIDs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(diffCommitIDs)

	pulls, err := fetchPullRequests(FetchPullRequestsLimitDefault)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(pulls))

	pullRequestTitle := fmt.Sprintf("Release %s", time.Now().Format("2006-01-02"))
	pullRequestBody := "# Pull Requests\n\n"
	for _, commitID := range diffCommitIDs {
		for _, pull := range pulls {
			if commitID == pull.GetMergeCommitSHA() {
				pullRequestBody += fmt.Sprintf("* %s (#%d) @%s\n", pull.GetTitle(), pull.GetNumber(), pull.GetUser().GetLogin())
			}
		}
	}

	fmt.Println(pullRequestBody)

	releasePullRequest, isCreated, err := createOrUpdatePullRequest(pullRequestTitle, pullRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	if isCreated {
		fmt.Printf("Created %s\n", releasePullRequest.GetHTMLURL())
	} else {
		fmt.Printf("Updated %s\n", releasePullRequest.GetHTMLURL())
	}
}

// Fetch the difference of commit IDs between develop and main
func fetchDiffCommitIDs() ([]string, error) {
	var diffCommitIDs []string

	comparison, _, err := githubClient.Repositories.CompareCommits(
		ctx,
		params.RepositoryOwner,
		params.RepositoryName,
		params.ProductionBranchName,
		params.DevelopmentBranchName,
	)
	if err != nil {
		return diffCommitIDs, err
	}

	for _, commit := range comparison.Commits {
		diffCommitIDs = append(diffCommitIDs, commit.GetSHA())
	}

	return diffCommitIDs, nil
}

// Fetch up to `limit` pull requests sorted by updated desc
func fetchPullRequests(limit int) ([]*github.PullRequest, error) {
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
func createOrUpdatePullRequest(title string, body string) (*github.PullRequest, bool, error) {
	var pullRequest *github.PullRequest
	isCreated := false

	releasePullRequests, _, err := githubClient.PullRequests.List(ctx, params.RepositoryOwner, params.RepositoryName, &github.PullRequestListOptions{
		Head:  params.DevelopmentBranchName,
		Base:  params.ProductionBranchName,
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
		Base:  &params.ProductionBranchName,
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
