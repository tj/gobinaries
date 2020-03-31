package semver

import (
	"errors"
	"strconv"
	"strings"
)

// ErrMalformedRange is returned when the format is incorrect.
var ErrMalformedRange = errors.New("malformed range format")

// Range is a set of operations used to match a version,
// currently only wildcards are supported.
type Range struct {
	Major string
	Minor string
	Patch string
}

// Match returns true if the range matches the version.
func (r Range) Match(v Version) bool {
	// major
	if r.Major != "*" {
		n, _ := strconv.ParseUint(r.Major, 10, 64)
		if v.Major != int(n) {
			return false
		}
	}

	// minor
	if r.Minor != "*" {
		n, _ := strconv.ParseUint(r.Minor, 10, 64)
		if v.Minor != int(n) {
			return false
		}
	}

	// patch
	if r.Patch != "*" {
		n, _ := strconv.ParseUint(r.Patch, 10, 64)
		if v.Patch != int(n) {
			return false
		}
	}

	return true
}

// ParseRange returns a parsed range.
func ParseRange(s string) (Range, error) {
	s = strings.TrimPrefix(s, "v")
	p := normalizeRangeParts(strings.Split(s, "."))

	return Range{
		Major: p[0],
		Minor: p[1],
		Patch: p[2],
	}, nil
}

// normalizeRangeParts returns normalized range parts,
// so that it is always a 3-tuple, with each version
// defaulting to "*".
func normalizeRangeParts(p []string) []string {
	// ensure always 3 parts
	for i := 0; i < 4-len(p); i++ {
		p = append(p, "")
	}

	// default to "*", convert "x" to "*"
	for i, s := range p {
		if s == "" || s == "x" {
			p[i] = "*"
		}
	}

	return p
}
