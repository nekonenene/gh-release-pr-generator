package cli

import (
	"fmt"
	"log"
	"time"
)

// ENTRY POINT of this package
func Exec() {
	ParseParameters()
	InitContextAndClient()

	diffCommitIDs, err := FetchDiffCommitIDs()
	if err != nil {
		log.Fatal(err)
	}

	pulls, err := FetchPullRequests(FetchPullRequestsLimitDefault)
	if err != nil {
		log.Fatal(err)
	}

	pullRequestTitle := fmt.Sprintf("Release %s", time.Now().Format("2006-01-02"))
	pullRequestBody := "# Pull Requests\n\n"
	for _, commitID := range diffCommitIDs {
		for _, pull := range pulls {
			if commitID == pull.GetMergeCommitSHA() {
				pullRequestBody += fmt.Sprintf("* %s (#%d) @%s\n", pull.GetTitle(), pull.GetNumber(), pull.GetUser().GetLogin())
			}
		}
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
