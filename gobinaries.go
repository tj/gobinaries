// Package gobinaries provides an HTTP server for on-demand Go binaries.
package gobinaries

import (
	"context"
	"errors"
	"io"
)

// ErrObjectNotFound is returned by Storage.Get() when no object is found for the specified key.
var ErrObjectNotFound = errors.New("no cloud storage object")

// ErrNoVersionMatch is returned by Resolver.Resolve() when no tag matches the requested version.
var ErrNoVersionMatch = errors.New("no matching version")

// ErrNoVersions is returned by Resolver.Resolve() when no versions are defined.
var ErrNoVersions = errors.New("no versions defined")

// Resolver is the interface used to resolver package versions.
type Resolver interface {
	Resolve(owner, repo, version string) (string, error)
}

// Storage is the interface used for storing compiled Go binaries.
type Storage interface {
	Create(context.Context, io.Reader, Binary) error
	Get(context.Context, Binary) (io.ReadCloser, error)
}

// Binary represents the details of a package binary.
type Binary struct {
	// Path is the command path such as "github.com/tj/staticgen/cmd/staticgen".
	Path string

	// Module path such as "github.com/tj/staticgen".
	Module string

	// Version is the version of the package.
	Version string

	// OS is the the target operating system.
	OS string

	// Arch is the target architecture.
	Arch string
}
