package semver_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/tj/gobinaries/semver"
)

// Test parsing of ranges.
func TestParseRange(t *testing.T) {
	t.Run("missing major", func(t *testing.T) {
		m, err := semver.ParseRange("")
		assert.NoError(t, err)
		assert.Equal(t, semver.Range{"*", "*", "*"}, m)
	})

	t.Run("missing minor", func(t *testing.T) {
		m, err := semver.ParseRange("1")
		assert.NoError(t, err)
		assert.Equal(t, semver.Range{"1", "*", "*"}, m)
	})

	t.Run("missing patch", func(t *testing.T) {
		m, err := semver.ParseRange("1.2")
		assert.NoError(t, err)
		assert.Equal(t, semver.Range{"1", "2", "*"}, m)
	})

	t.Run("exact match", func(t *testing.T) {
		m, err := semver.ParseRange("1.2.3")
		assert.NoError(t, err)
		assert.Equal(t, semver.Range{"1", "2", "3"}, m)
	})

	t.Run("globs", func(t *testing.T) {
		m, err := semver.ParseRange("1.*")
		assert.NoError(t, err)
		assert.Equal(t, semver.Range{"1", "*", "*"}, m)
	})

	t.Run("wildcard", func(t *testing.T) {
		m, err := semver.ParseRange("1.x")
		assert.NoError(t, err)
		assert.Equal(t, semver.Range{"1", "*", "*"}, m)
	})
}

// Test range matching.
func TestRange_Match(t *testing.T) {
	t.Run("exact match", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("2.0.5")
		assert.NoError(t, err)
		assert.True(t, m.Match(v))
	})

	t.Run("exact match with leading v", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("v2.0.5")
		assert.NoError(t, err)
		assert.True(t, m.Match(v))
	})

	t.Run("exact mismatch", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("2.0.2")
		assert.NoError(t, err)
		assert.False(t, m.Match(v))
	})

	t.Run("major minor mismatch", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("3.0")
		assert.NoError(t, err)
		assert.False(t, m.Match(v))
	})

	t.Run("wildcard minor mismatch", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("3.x")
		assert.NoError(t, err)
		assert.False(t, m.Match(v))
	})

	t.Run("wildcard patch mismatch", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("2.1.x")
		assert.NoError(t, err)
		assert.False(t, m.Match(v))
	})

	t.Run("major minor match", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("2.0")
		assert.NoError(t, err)
		assert.True(t, m.Match(v))
	})

	t.Run("wildcard minor match", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("2.x")
		assert.NoError(t, err)
		assert.True(t, m.Match(v))
	})

	t.Run("wildcard patch match", func(t *testing.T) {
		v := semver.Version{Major: 2, Minor: 0, Patch: 5}
		m, err := semver.ParseRange("2.0.x")
		assert.NoError(t, err)
		assert.True(t, m.Match(v))
	})
}
