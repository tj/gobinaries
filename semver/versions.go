package semver

// Versions is a slice of versions.
type Versions []Version

func (a Versions) Len() int           { return len(a) }
func (a Versions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Versions) Less(i, j int) bool { return a[i].Compare(a[j]) == -1 }
