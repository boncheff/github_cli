package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github_cli/types"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// IssuesURL ...
const IssuesURL = "https://api.github.com/search/issues"

// SearchIssues ...
func SearchIssues(repo, open *string, terms []string) {

	terms = append(terms, "repo:"+*repo)
	terms = append(terms, *open)
	q := url.QueryEscape(strings.Join(terms, " "))

	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatalf("search query failed: %s", resp.Status)
	}

	var result types.IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	fmt.Printf("%d issue(s) for repo %q with state %q:\n\n", result.TotalCount, *repo, *open)
	for _, item := range result.Items {
		issueURL := fmt.Sprintf("https://github.com/%s/pull/%d", *repo, item.Number)
		title := formatTitle(item.Title)

		fmt.Printf(
			"#%-10d %-25s %-85s %-12s %-50s\n",
			item.Number, item.User.Login, title, time.Time.Format(item.CreatedAt, "02-01-2006"), issueURL,
		)
	}
}

// UpdateIssue ...
func UpdateIssue(repo, state, issueNumber *string) {
	patchURL := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s", *repo, *issueNumber)

	var jsonStr = []byte(`{"state": "` + *state + `" }`)

	auth := getAuth()
	req, err := http.NewRequest("PATCH", patchURL, bytes.NewBuffer(jsonStr))
	req.SetBasicAuth(auth.Username, auth.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	if strings.Contains(s, "Not Found") {
		log.Fatalf(" > Issue %s not found", *issueNumber)
	}
	fmt.Println(" > OK")

}

func formatTitle(itemTitle string) string {
	var title string
	if len(itemTitle) > 80 {
		title = itemTitle[:80]
	} else {
		title = itemTitle
	}

	return title
}

func getAuth() types.GithubAuth {
	yamlFile, err := ioutil.ReadFile(".github.yaml")
	if err != nil {
		log.Fatalf("Error reading github.yaml file: %s", err)
	}

	var auth types.GithubAuth
	err = yaml.Unmarshal(yamlFile, &auth)
	if err != nil {
		log.Fatalf("Unmarshal error: %s", err)
	}

	return auth
}
