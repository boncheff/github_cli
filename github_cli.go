package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github_cli/types"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

func main() {
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	repo := listCommand.String("repo", "golang/go", "the repo for which to search for issues")
	open := listCommand.String("open", "is:open", "search only for open issues")
	terms := listCommand.String("terms", "", "search terms to include when querying issues")

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	state := updateCommand.String("state", "", "the new state of the ticket")
	issueNumber := updateCommand.String("issue", "", "the issue to update")
	updateRepo := updateCommand.String("repo", "golang/go", "the repo for which to update fields")

	if len(os.Args) == 1 {
		fmt.Println("usage: github_cli <command> [<args>]")
		fmt.Println("Supported commands: ")
		fmt.Println(" list   List all issues")
		fmt.Println(" update Update an issue")
		return
	}

	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
		queryTerms := strings.Split(*terms, " ")
		SearchIssues(repo, open, queryTerms)

	case "update":
		updateCommand.Parse(os.Args[2:])
		UpdateIssue(updateRepo, state, issueNumber)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

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

func formatTitle(itemTitle string) string {
	var title string
	if len(itemTitle) > 80 {
		title = itemTitle[:80]
	} else {
		title = itemTitle
	}

	return title
}
