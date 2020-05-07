package server

import (
	"fmt"
	"strconv"
	"strings"
)

// getMajorVersion tries to detect the major version of the package.
func getMajorVersion(tag string) (int, error) {
	major := strings.Split(tag, ".")[0]
	if len(major) < 1 {
		return 0, fmt.Errorf("invalid major version")
	}
	return strconv.Atoi(major[1:])
}

// parsePackage returns package information parsed from the path.
func parsePackage(path string) (pkg, mod, version, bin string) {
	p := strings.Split(path, "@")
	version = "master"

	// pkg
	pkg = normalizePackage(p[0])

	// mod
	modp := strings.Split(pkg, "/")
	if len(modp) >= 3 {
		mod = strings.Join(modp[:3], "/")
	}

	// version after @
	if len(p) > 1 {
		version = p[1]
	}

	// binary name from pkg
	p = strings.Split(pkg, "/")
	bin = p[len(p)-1]
	return
}

// normalizePackage returns a normalized package, where "https://github.com/" is implied.
func normalizePackage(pkg string) string {
	// ignore leading https://
	pkg = strings.Replace(pkg, "https://", "", 1)

	// ignore leading github.com/
	pkg = strings.Replace(pkg, "github.com/", "", 1)

	// implicit github.com
	pkg = "github.com/" + pkg

	return pkg
}
