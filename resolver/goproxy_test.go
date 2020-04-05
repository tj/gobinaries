package resolver_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/tj/go/env"
	"github.com/tj/gobinaries"
	"github.com/tj/gobinaries/resolver"
)

// newGoProxyResolver returns a new GitHub resolver.
func newGoProxyResolver() gobinaries.Resolver {
	return &resolver.GoProxy{
		URL: env.Get("GOPROXY"),
	}
}

// Test resolver.
func TestGoProxy_Resolve(t *testing.T) {
	r := newGoProxyResolver()

	t.Run("exact match", func(t *testing.T) {
		repo := resolver.Repository{
			Location: "github.com",
			Owner:    "tj",
			Project:  "d3-bar",
			Version:  "v1.8.0",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("exact match without leading v", func(t *testing.T) {
		repo := resolver.Repository{
			Location: "github.com",
			Owner:    "tj",
			Project:  "d3-bar",
			Version:  "1.8.0",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("major wildcard match", func(t *testing.T) {
		repo := resolver.Repository{
			Location: "github.com",
			Owner:    "tj",
			Project:  "d3-bar",
			Version:  "1.x",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.8.0", v)
	})

	t.Run("minor wildcard match", func(t *testing.T) {
		repo := resolver.Repository{
			Location: "github.com",
			Owner:    "tj",
			Project:  "d3-bar",
			Version:  "1.6.x",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.6.0", v)
	})

	t.Run("minor match", func(t *testing.T) {
		repo := resolver.Repository{
			Location: "github.com",
			Owner:    "tj",
			Project:  "d3-bar",
			Version:  "1.6",
		}
		v, err := r.Resolve(repo)
		assert.NoError(t, err)
		assert.Equal(t, "v1.6.0", v)
	})
}
