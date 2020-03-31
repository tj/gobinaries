package resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v28/github"

	"github.com/tj/gobinaries"
	"github.com/tj/gobinaries/semver"
)

// GitHub is an implementation of Versioner for
// listing GitHub repository tags defined for a project.
type GitHub struct {
	// Client is the GitHub client.
	Client *github.Client
}

// Resolve implementation.
func (g *GitHub) Resolve(owner, repo, version string) (string, error) {
	// fetch tags
	tags, err := g.versions(owner, repo)
	if err != nil {
		return "", err
	}

	// convert to semver, ignoring malformed
	var versions []semver.Version
	for _, t := range tags {
		if v, err := semver.Parse(t); err == nil {
			versions = append(versions, v)
		}
	}

	// master special-case
	if version == "master" {
		return versions[0].String(), nil
	}

	// match requested semver range
	vr, err := semver.ParseRange(version)
	if err != nil {
		return "", fmt.Errorf("parsing version range: %w", err)
	}

	for _, v := range versions {
		if vr.Match(v) {
			return v.String(), nil
		}
	}

	return "", gobinaries.ErrNoVersionMatch
}

// versions returns the versions of a repository.
func (g *GitHub) versions(owner, repo string) (versions []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	page := 1

	for {
		options := &github.ListOptions{
			Page:    page,
			PerPage: 100,
		}

		tags, _, err := g.Client.Repositories.ListTags(ctx, owner, repo, options)
		if err != nil {
			return nil, fmt.Errorf("listing tags: %w", err)
		}

		if len(tags) == 0 {
			break
		}

		for _, t := range tags {
			versions = append(versions, t.GetName())
		}

		page++
	}

	if len(versions) == 0 {
		return nil, gobinaries.ErrNoVersions
	}

	return
}
