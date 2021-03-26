package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

func Exec() {
	ParseParameters()

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.GitHubAPIToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)
	githubClient := github.NewClient(httpClient)

	// Fetch the difference of commit IDs between develop and main
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

	pulls, _, err := githubClient.PullRequests.List(ctx, params.RepositoryOwner, params.RepositoryName, &github.PullRequestListOptions{
		Base:      params.DevelopmentBranchName,
		State:     "closed",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

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
