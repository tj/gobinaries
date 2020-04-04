// Package build provides Go package building.
package build

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tj/gobinaries"
)

// environMap returns a map of environment variables.
var environMap map[string]string

// init initializes the environment variable map.
func init() {
	environMap = make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		environMap[parts[0]] = parts[1]
	}
}

// environWhitelist is a list of environment variables to include in Go sub-commands.
var environWhitelist = []string{
	"PATH",
	"HOME",
	"PWD",
	"GOPATH",
	"GOLANG_VERSION",
	"TMPDIR",
}

// ErrNotExecutable is returned when the package path provided does not produce a binary.
var ErrNotExecutable = errors.New("not executable")

// Error represents a build error.
type Error struct {
	err    error
	stderr string
}

// Error implementation.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.err.Error(), e.stderr)
}

// Write a package binary to w.
func Write(w io.Writer, bin gobinaries.Binary) error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("user home dir: %w", err)
	}

	// remove the old go.mod if there is one
	err = os.Remove(filepath.Join(dir, "go.mod"))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing go.mod: %w", err)
	}

	// create a go.mod file, this is currently required
	// in order to install a package with a specified version
	err = addModule(dir)
	if err != nil {
		return fmt.Errorf("initializing module: %w", err)
	}

	// add the dependency
	err = addModuleDep(dir, bin.Module+"@"+bin.Version)
	if err != nil {
		return fmt.Errorf("adding dependency: %w", err)
	}

	// tmpfile for the binary
	dst, err := tempFilename()
	if err != nil {
		return fmt.Errorf("creating tempfile: %w", err)
	}

	// build the binary
	err = buildBinary(dir, dst, bin)
	if err != nil {
		return fmt.Errorf("building: %w", err)
	}

	// check permissions and copy it to w
	f, err := os.Open(dst)
	if err != nil {
		return fmt.Errorf("opening: %w", err)
	}

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stating: %w", err)
	}

	if !isExecutable(info.Mode()) {
		return ErrNotExecutable
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return fmt.Errorf("copying: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("closing: %w", err)
	}

	err = os.Remove(dst)
	if err != nil {
		return fmt.Errorf("removing tempfile: %w", err)
	}

	return nil
}

// ClearCache removes the module cache.
func ClearCache() error {
	cmd := exec.Command("go", "clean", "--modcache")
	return cmd.Run()
}

// isExecutable returns true if the exec bit is set for u/g/o.
func isExecutable(mode os.FileMode) bool {
	return mode.Perm()&0111 == 0111
}

// addModule initializes a new go module in the given dir. This is apparently
// necessary to build using Go modules since `go build` does not support
// semver, awkward UX but oh well.
func addModule(dir string) error {
	cmd := exec.Command("go", "mod", "init", "github.com/gobinary")
	cmd.Env = environ()
	cmd.Env = append(cmd.Env, "GO111MODULE=on")
	cmd.Dir = dir
	return command(cmd)
}

// addModuleDep creates a module dependency.
func addModuleDep(dir, dep string) error {
	cmd := exec.Command("go", "mod", "edit", "-require", dep)
	cmd.Env = environ()
	cmd.Env = append(cmd.Env, "GO111MODULE=on")
	cmd.Dir = dir
	return command(cmd)
}

// buildBinary performs a `go build` and outputs the binary to dst.
func buildBinary(dir, dst string, bin gobinaries.Binary) error {
	ldflags := fmt.Sprintf("-X main.version=%s", bin.Version)
	cmd := exec.Command("go", "build", "-o", dst, "-ldflags", ldflags, bin.Path)
	cmd.Env = environ()
	cmd.Env = append(cmd.Env, "GO111MODULE=on")
	cmd.Env = append(cmd.Env, "GOOS="+bin.OS)
	cmd.Env = append(cmd.Env, "GOARCH="+bin.Arch)
	cmd.Dir = dir
	return command(cmd)
}

// command executes a command and capture stderr.
func command(cmd *exec.Cmd) error {
	var w strings.Builder
	cmd.Stderr = &w
	err := cmd.Run()
	if err != nil {
		return Error{
			err:    err,
			stderr: strings.TrimSpace(w.String()),
		}
	}
	return nil
}

// tempFilename returns a new temporary file name.
func tempFilename() (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "gobinary")
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	return f.Name(), nil
}

// environ returns the environment variables for Go sub-commands.
func environ() (env []string) {
	for _, name := range environWhitelist {
		env = append(env, name+"="+environMap[name])
	}
	return
}
