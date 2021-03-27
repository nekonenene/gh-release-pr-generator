package cmd

import (
	"fmt"
	"log"

	"github.com/google/go-github/v34/github"
)

// ENTRY POINT of this package
func Exec() {
	ParseParameters()

	err := InitContextAndGitHubClient()
	if err != nil {
		log.Fatal(err)
	}

	diffCommitIDs, err := FetchDiffCommitIDs()
	if err != nil {
		log.Fatal(err)
	}

	// Exit because GitHub API would refuse to create new pull request when no differences between branches
	if len(diffCommitIDs) == 0 {
		fmt.Printf("No differences between %s and %s branches\n", params.DevelopmentBranchName, params.ProductionBranchName)
		return
	}

	closedPulls, err := FetchClosedPullRequests(params.FetchPullRequestsLimit)
	if err != nil {
		log.Fatal(err)
	}

	// Select pull requests which has not yet been merged into the main branch from closed ones
	var taegtPulls []*github.PullRequest
	for _, commitID := range diffCommitIDs {
		for _, pull := range closedPulls {
			if commitID == pull.GetMergeCommitSHA() {
				taegtPulls = append(taegtPulls, pull)
			}
		}
	}

	pullRequestTitle, pullRequestBody, err := ConstructTitleAndBody(taegtPulls)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[Title]")
	fmt.Println(pullRequestTitle)
	fmt.Println()
	fmt.Println("[Body]")
	fmt.Println(pullRequestBody)
	fmt.Println()

	releasePullRequest, isCreated, err := CreateOrUpdatePullRequest(pullRequestTitle, pullRequestBody)
	if err != nil {
		log.Fatal(err)
	}

	if isCreated {
		fmt.Printf("Created %s\n", releasePullRequest.GetHTMLURL())
	} else {
		fmt.Printf("Updated %s\n", releasePullRequest.GetHTMLURL())
	}
}
