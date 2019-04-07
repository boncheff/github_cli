package main

import (
	"flag"
	"fmt"
	"github_cli/handlers"
	"os"
	"strings"
)

func main() {
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	repo := listCommand.String("repo", "golang/go", "the repo for which to search for issues")
	listState := listCommand.String("state", "is:open", "search only for issues with given state")
	terms := listCommand.String("terms", "", "search terms to include when querying issues")

	updateCommand := flag.NewFlagSet("update", flag.ExitOnError)
	state := updateCommand.String("state", "", "the new state of the ticket")
	issueNumber := updateCommand.String("issue", "", "the issue to update")
	updateRepo := updateCommand.String("repo", "golang/go", "the repo for which to update fields")

	createCommand := flag.NewFlagSet("create", flag.ExitOnError)
	createTitle := createCommand.String("title", "", "the title of the issue")
	createBody := createCommand.String("body", "", "the body of the issue")
	createRepo := createCommand.String("repo", "golang/go", "the repo for which to create a new issue")

	if len(os.Args) == 1 {
		fmt.Println("usage: github_cli <command> [<args>]")
		fmt.Println("Supported commands: ")
		fmt.Println(" list   List all issues")
		fmt.Println(" update Update an issue")
		fmt.Println(" create Create an issue")
		return
	}

	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
		queryTerms := strings.Split(*terms, " ")
		handlers.SearchIssues(repo, listState, queryTerms)

	case "update":
		updateCommand.Parse(os.Args[2:])
		handlers.UpdateIssue(updateRepo, state, issueNumber)

	case "create":
		createCommand.Parse(os.Args[2:])
		handlers.CreateIssue(createRepo, createTitle, createBody)

	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

}
