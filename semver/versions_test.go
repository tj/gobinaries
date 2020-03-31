package semver_test

import (
	"sort"
	"testing"

	"github.com/tj/assert"

	"github.com/tj/gobinaries/semver"
)

// Test sorting.
func TestSort(t *testing.T) {
	versions := semver.Versions{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 1, Minor: 0, Patch: 5},
		{Major: 1, Minor: 2, Patch: 0},
		{Major: 1, Minor: 2, Patch: 3},
		{Major: 3, Minor: 0, Patch: 0},
		{Major: 2, Minor: 0, Patch: 1},
		{Major: 2, Minor: 0, Patch: 0},
		{Major: 2, Minor: 1, Patch: 0},
	}

	sort.Sort(sort.Reverse(versions))

	first := versions[0]
	assert.Equal(t, semver.Version{Major: 3, Minor: 0, Patch: 0}, first)

	last := versions[len(versions)-1]
	assert.Equal(t, semver.Version{Major: 1, Minor: 0, Patch: 0}, last)
}
