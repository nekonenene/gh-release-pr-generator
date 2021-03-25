package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func Exec() {
	ParseParameters()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.GitHubAPIToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	pulls, _, err := client.PullRequests.List(ctx, params.RepositoryOwner, params.RepositoryName, &github.PullRequestListOptions{
		State: "closed",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created new repo: %v\n", pulls)
}
