package resolver

import (
	"bufio"
	"fmt"
	"net/http"
	"sort"

	"github.com/tj/gobinaries/semver"
)

// modListURLTemplate is the well known location for getting a go modules
// version list from a GOPROXY.
const modListURLTemplate = "%s/%s/%s/%s/@v/list"

// GoProxy holds the location of the GOPROXY a user wants to resolve
// versions with.
type GoProxy struct {
	URL string
}

// Resolve attempts to resolve a properly formatted repository / verison
// via the GOPROXY (acting like the request is a module).
func (g *GoProxy) Resolve(repo Repository) (string, error) {
	repoURL := fmt.Sprintf(modListURLTemplate, g.URL, repo.Location, repo.Owner, repo.Project)

	resp, err := http.Get(repoURL)
	if err != nil {
		return "", fmt.Errorf("goproxy failed to get verions: %w", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var versions []semver.Version
	for scanner.Scan() {
		token := string(scanner.Bytes())
		if v, err := semver.Parse(token); err == nil {
			versions = append(versions, v)
		}
	}

	// sort ascending for all versions
	sort.Sort(sort.Reverse(semver.Versions(versions)))

	vr, err := semver.ParseRange(repo.Version)
	if err != nil {
		return "", fmt.Errorf("parsing version range: %w", err)
	}

	for _, v := range versions {
		if vr.Match(v) {
			return v.String(), nil
		}
	}

	return "", ErrNoVersionMatch
}
