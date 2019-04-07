package clients

import (
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

// GithubClient ...
type GithubClient struct {
	*BaseClient
}

// IssuesURL ...
const IssuesURL = "https://api.github.com/search/issues"

// SearchIssues ...
func (h *GithubClient) SearchIssues(repo, open *string, terms []string) {

	terms = append(terms, "repo:"+*repo)
	terms = append(terms, *open)
	q := url.QueryEscape(strings.Join(terms, " "))

	var result types.IssuesSearchResult
	resp, err := h.Req.
		SetResult(&result).
		Get(IssuesURL + "?q=" + q)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode() != http.StatusOK {
		log.Fatalf("search query failed: %s", resp.Error())
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
func (h *GithubClient) UpdateIssue(repo, state, issueNumber *string) {
	patchURL := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s", *repo, *issueNumber)

	var jsonStr = []byte(`{"state": "` + *state + `" }`)

	auth := getAuth()
	req, err := h.Req.
		SetBody(jsonStr).
		SetBasicAuth(auth.Username, auth.Password).
		SetHeader("Content-Type", "application/json").
		Patch(patchURL)

	if err != nil {
		log.Fatal(err)
	}

	if req.StatusCode() == http.StatusNotFound {
		log.Fatalf(" > Issue %s not found", *issueNumber)
	}
	fmt.Println(" > OK")

}

// CreateIssue ...
func (h *GithubClient) CreateIssue(repo, title, body *string) {
	postURL := fmt.Sprintf("https://api.github.com/repos/%s/issues", *repo)

	var jsonStr = []byte(`{"title": "` + *title + `", "body": "` + *body + `"}`)

	var postResult struct {
		Number int `json:"number"`
	}

	auth := getAuth()
	req, err := h.Req.
		SetBody(jsonStr).
		SetBasicAuth(auth.Username, auth.Password).
		SetResult(&postResult).
		SetHeader("Content-Type", "application/json").
		Post(postURL)

	if err != nil {
		log.Fatal(err)
	}

	if req.StatusCode() != http.StatusCreated {
		log.Fatalf("Error creating issue")
	}

	fmt.Printf(" > Created issue #%d\n", postResult.Number)

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
