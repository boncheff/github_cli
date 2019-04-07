package clients

import (
	resty "gopkg.in/resty.v0"
)

// Client ...
type Client interface {
	SearchIssues(*string, *string, []string)
	UpdateIssue(*string, *string, *string)
	CreateIssue(*string, *string, *string)
}

// BaseClient ...
type BaseClient struct {
	Req *resty.Request
}
