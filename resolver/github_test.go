package resolver_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/tj/assert"
	"github.com/tj/go/env"
	"golang.org/x/oauth2"

	"github.com/tj/gobinaries"
	"github.com/tj/gobinaries/resolver"
)

// newResolver returns a new GitHub resolver.
func newResolver() gobinaries.Resolver {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: env.Get("GITHUB_TOKEN"),
		},
	)

	return &resolver.GitHub{
		Client: github.NewClient(oauth2.NewClient(ctx, ts)),
	}
}

// Test resolver.
func TestGitHub_Resolve(t *testing.T) {
	r := newResolver()

	t.Run("exact match", func(t *testing.T) {
		v, err := r.Resolve("tj", "d3-bar", "v1.8.0")
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("exact match without leading v", func(t *testing.T) {
		v, err := r.Resolve("tj", "d3-bar", "1.8.0")
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("major wildcard match", func(t *testing.T) {
		v, err := r.Resolve("tj", "d3-bar", "1.x")
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("minor wildcard match", func(t *testing.T) {
		v, err := r.Resolve("tj", "d3-bar", "1.6.x")
		assert.NoError(t, err)
		assert.Equal(t, "v1.6.0", v)
	})

	t.Run("minor match", func(t *testing.T) {
		v, err := r.Resolve("tj", "d3-bar", "1.6")
		assert.NoError(t, err)
		assert.Equal(t, "v1.6.0", v)
	})

	t.Run("master", func(t *testing.T) {
		v, err := r.Resolve("tj", "d3-bar", "master")
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})
}
