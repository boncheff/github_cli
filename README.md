# Github CLI Interface

This is a simple CLI written in [Go](https://golang.org/) that interacts with Github using the [Github's issues CRUD API](https://developer.github.com/v3/issues/)

## Prerequisites

In order to be able to run the program, you need to create a *.github.yaml* file in the current directory. This file is used to authenticate the user when talking to GitHub and has the following format:

```
username: <github username>

password: <github password>
```

## Getting Started

Currently the CLI suports the following actions:
1. Creating new issues
2. Listing issues by repo, state and terms
3. Updating an issue's state (open/closed) by issue number and repo

### Examples

```
go build github_cli.go

./github_cli list --repo golang/go --state is:open

./github_cli create --repo golang/go --title "Issue Title" --body "Issue body"

./github_cli.go update --repo golang/go --issue 359 --state closed
```
