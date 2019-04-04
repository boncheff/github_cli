# Github CLI Interface

This is a simple CLI written in Go that interacts with Github which (will eventually)  make full use of [Github's issues CRUD API](https://developer.github.com/v3/issues/)

It requires that a .github.yaml file be present in the current directory. This file is used to authenticate the user and has the following format:

> username: github username

> password: github password

Currently the only supported actions are:
1. Listing issues by repo, state and terms
2. Updating an issue's state (open/closed) by issue number and repo

## How to use

> go build github_cli.go
> ./github_cli list --repo golang/go --open is:open
