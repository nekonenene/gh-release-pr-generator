package cli

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func Exec() {
	ParseParameters()

	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.GitHubAPIToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)
	githubClient := github.NewClient(httpClient)

	pulls, _, err := githubClient.PullRequests.List(ctx, params.RepositoryOwner, params.RepositoryName, &github.PullRequestListOptions{
		State: "closed",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, pull := range pulls {
		fmt.Println(pull.GetTitle())
		fmt.Println(pull.GetUser().GetLogin())

		// マージされずに閉じられたものはスキップ
		if (pull.GetMergedAt() == time.Time{}) {
			continue
		}
		fmt.Println(pull.GetMergedAt())
	}
}
