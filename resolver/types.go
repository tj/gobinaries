package resolver

import "errors"

// ErrNoVersionMatch is returned by Resolver.Resolve() when no tag matches the requested version.
var ErrNoVersionMatch = errors.New("no matching version")

// ErrNoVersions is returned by Resolver.Resolve() when no versions are defined.
var ErrNoVersions = errors.New("no versions defined")

// Respository holds the details of where,
// who, and what version of a repository we
// are trying to resolve. Location is not used
// by GitHub as we KNOW the location is always
// GitHub.
type Repository struct {
	Location string
	Owner    string
	Project  string
	Version  string
}
