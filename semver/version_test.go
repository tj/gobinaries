package semver_test

import (
	"testing"

	"github.com/tj/assert"

	"github.com/tj/gobinaries/semver"
)

// Test parsing of versions.
func TestParse(t *testing.T) {
	t.Run("missing major", func(t *testing.T) {
		_, err := semver.Parse("")
		assert.Equal(t, err, semver.ErrMalformed)
	})

	t.Run("missing minor", func(t *testing.T) {
		_, err := semver.Parse("1")
		assert.Equal(t, err, semver.ErrMalformed)
	})

	t.Run("missing patch", func(t *testing.T) {
		v, err := semver.Parse("1.2")
		assert.NoError(t, err)
		assert.Equal(t, semver.Version{Major: 1, Minor: 2, Patch: 0, Input: "1.2"}, v)
	})

	t.Run("valid", func(t *testing.T) {
		v, err := semver.Parse("1.2.3")
		assert.NoError(t, err)
		assert.Equal(t, semver.Version{Major: 1, Minor: 2, Patch: 3, Input: "1.2.3"}, v)
	})

	t.Run("valid with leading v", func(t *testing.T) {
		v, err := semver.Parse("v1.2.3")
		assert.NoError(t, err)
		assert.Equal(t, semver.Version{Major: 1, Minor: 2, Patch: 3, Input: "v1.2.3"}, v)
		assert.Equal(t, "v1.2.3", v.String())
	})
}
