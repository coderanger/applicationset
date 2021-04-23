package repo_host

import (
	"context"
	"regexp"
)

// An abstract repository from an API provider.
type HostedRepo struct {
	Organization string
	Repository   string
	URL          string
	Branch       string
	Labels       []string
}

type RepoHostService interface {
	ListRepos(context.Context) ([]*HostedRepo, error)
	RepoHasPath(context.Context, *HostedRepo, string) (bool, error)
}

// A compiled version of RepoHostGeneratorFilter for performance.
type Filter struct {
	RepositoryMatch *regexp.Regexp
	PathExists      *string
	LabelMatch      *regexp.Regexp
}
