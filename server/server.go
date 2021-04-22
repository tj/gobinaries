// Package server provides an HTTP server for on-demand Go binaries.
package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/apex/log"
	"github.com/tj/go/http/request"
	"github.com/tj/go/http/response"

	"github.com/tj/gobinaries"
	"github.com/tj/gobinaries/build"
)

// Server is the binary server.
type Server struct {
	// URL is the API endpoint URL.
	URL string

	// Static file directory.
	Static string

	// Store is the object storage.
	Storage gobinaries.Storage

	// Resolver is the version resolver.
	Resolver gobinaries.Resolver

	once      sync.Once
	templates *template.Template
}

// ServeHTTP implementation.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.once.Do(func() {
		s.templates = template.Must(template.ParseGlob("templates/*"))
	})

	path := r.URL.Path

	// invalid method
	if r.Method != "GET" {
		response.MethodNotAllowed(w)
		return
	}

	// normalize index
	if path == "/" {
		path = "/index.html"
	}

	// health check
	if path == "/_health" {
		response.OK(w, ":)")
		return
	}

	// serve binary
	if strings.HasPrefix(path, "/binary/") {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/binary/")
		s.getBinary(w, r)
		return
	}

	// check if we have a static file,
	// serve it before we try to fetch
	// information from Github
	file := filepath.Join(s.Static, path)
	info, err := os.Stat(file)
	if err == nil && info.Mode().IsRegular() {
		http.ServeFile(w, r, file)
		return
	}

	// serve installation script
	s.getScript(w, r)
}

// getScript takes a package path such as "tj/staticgen"
// or "github.com/tj/staticgen@1.x" and resolves the requested
// version, responding with an installation script to request
// the binary built for the user's machine.
//
// Known errors respond with shell scripts as well,
// in order to provide nicer in-shell error messages,
// otherwise the curl request will silently fail.
func (s *Server) getScript(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	pkg, mod, version, bin := parsePackage(path)

	if pkg == "" {
		response.BadRequest(w)
		return
	}

	parts := strings.Split(pkg, "/")

	if len(parts) < 3 {
		response.BadRequest(w)
		return
	}

	owner := parts[1]
	repo := parts[2]

	logs := log.WithFields(log.Fields{
		"ip":      r.Header.Get("CF-Connecting-IP"),
		"package": pkg,
		"module":  mod,
		"owner":   owner,
		"repo":    repo,
		"binary":  bin,
		"version": version,
	})

	logs.Info("resolving version")
	resolved, err := s.Resolver.Resolve(owner, repo, version)

	if err == gobinaries.ErrNoVersions {
		logs.Warn("no tags")
		s.render(w, "error.sh", "Repository has no tags")
		return
	}

	if err == gobinaries.ErrNoVersionMatch {
		logs.Warn("no match")
		s.render(w, "error.sh", "Repository has no tags matching the requested version")
		return
	}

	if err != nil {
		logs.WithError(err).Error("error resolving")
		s.render(w, "error.sh", "Failed to resolve requested version")
		return
	}

	logs = logs.WithField("resolved", resolved)
	logs.Info("resolved version")

	// rename package into go mod compatible name if v2 and above
	major, err := getMajorVersion(resolved)
	if err == nil && major > 1 {
		modp := strings.Split(pkg, "/")
		if len(modp) >= 3 {
			mod := strings.Join(modp[:3], "/")
			nested := strings.Join(modp[3:], "/")
			pkg = fmt.Sprintf("%s/v%d/%s", mod, major, nested)
		}
	}

	templateData := struct {
		URL             string
		Package         string
		Binary          string
		OriginalVersion string
		Version         string
	}{
		URL:             s.URL,
		Package:         pkg,
		Binary:          bin,
		OriginalVersion: version,
		Version:         resolved,
	}
	useragent := r.Header.Get("User-Agent")
	ispowershell := strings.Contains(strings.ToUpper(useragent), strings.ToUpper("POWERSHELL"))
	iswindows := strings.Contains(strings.ToUpper(useragent), strings.ToUpper("WINDOWS"))

	if ispowershell && iswindows {
		s.render(w, "install.ps1", templateData)
	} else {
		s.render(w, "install.sh", templateData)
	}
}

// getBinary builds and responds with the requested package binary,
// with the following required query-string parameters:
//
// - os
// - arch
// - version
//
// For example "github.com/tj/triage/cmd/triage?os=linux&arch=amd64&version=1.0.0".
//
func (s *Server) getBinary(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	pkg := strings.TrimPrefix(r.URL.Path, "/")

	if pkg == "" {
		response.BadRequest(w)
		return
	}

	goos := request.Param(r, "os")
	if goos == "" {
		response.BadRequest(w, "`os` parameter required")
		return
	}

	arch := request.Param(r, "arch")
	if arch == "" {
		response.BadRequest(w, "`arch` parameter required")
		return
	}

	version := request.Param(r, "version")
	if version == "" {
		response.BadRequest(w, "`version` parameter required")
		return
	}

	_, mod, _, _ := parsePackage(pkg)
	logs := log.WithFields(log.Fields{
		"ip":      r.Header.Get("CF-Connecting-IP"),
		"package": pkg,
		"module":  mod,
		"os":      goos,
		"arch":    arch,
		"version": version,
	})

	bin := gobinaries.Binary{
		Path:    pkg,
		Module:  mod,
		Version: version,
		OS:      goos,
		Arch:    arch,
	}

	// respond with the object if it already exists in storage
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	obj, err := s.Storage.Get(ctx, bin)
	if err == nil {
		logs.Info("serving from storage")
		immutable(w)
		_, _ = io.Copy(w, obj)
		return
	}

	// build the binary, writing it to the response
	// and buffering for cloud storage
	var buf bytes.Buffer
	logs.Info("building package")
	immutable(w)
	err = build.Write(io.MultiWriter(w, &buf), bin)
	if err != nil {
		logs.WithError(err).Error("building")
		response.InternalServerError(w)
		return
	}
	logs.WithField("duration", duration(start)).Info("built package")

	// store the binary
	start = time.Now()
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	logs.Info("storing package")
	err = s.Storage.Create(ctx, &buf, bin)
	if err == nil {
		logs.WithField("duration", duration(start)).Info("stored package")
	} else {
		logs.WithError(err).Error("storing binary")
	}

	// clear module cache
	start = time.Now()
	err = build.ClearCache()
	if err == nil {
		logs.WithField("duration", duration(start)).Info("cleared cache")
	} else {
		logs.WithError(err).Error("clearing cache")
	}
}

// render template.
func (s *Server) render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "application/x-sh")
	w.Header().Set("Cache-Control", "no-store")
	s.templates.ExecuteTemplate(w, name, data)
}

// immutable sets immutability header fields.
func immutable(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "max-age=31536000, immutable")
}

// duration returns the duration since start in milliseconds.
func duration(start time.Time) int {
	return int(time.Since(start) / time.Millisecond)
}
