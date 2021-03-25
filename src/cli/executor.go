package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Exec() {
	ParseParameters()

	url := apiHost + "/repos/" + params.RepositoryOwner + "/" + params.RepositoryName + "/pulls?state=closed"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", params.GitHubAPIToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch pull requests: %v", err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bodyBytes))
}
