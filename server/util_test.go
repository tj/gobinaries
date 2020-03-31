package server

import (
	"testing"

	"github.com/tj/assert"
)

// Test parsing package paths with top-level commands.
func TestParsePackage_root(t *testing.T) {
	t.Run("with version", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("https://github.com/tj/letterbox@v1.2.0")
		assert.Equal(t, "github.com/tj/letterbox", pkg)
		assert.Equal(t, "github.com/tj/letterbox", mod)
		assert.Equal(t, "v1.2.0", version)
		assert.Equal(t, "letterbox", bin)
	})

	t.Run("with https://github.com", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("https://github.com/tj/letterbox")
		assert.Equal(t, "github.com/tj/letterbox", pkg)
		assert.Equal(t, "github.com/tj/letterbox", mod)
		assert.Equal(t, "master", version)
		assert.Equal(t, "letterbox", bin)
	})

	t.Run("with github.com", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("github.com/tj/letterbox")
		assert.Equal(t, "github.com/tj/letterbox", pkg)
		assert.Equal(t, "github.com/tj/letterbox", mod)
		assert.Equal(t, "master", version)
		assert.Equal(t, "letterbox", bin)
	})

	t.Run("without host", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("tj/letterbox")
		assert.Equal(t, "github.com/tj/letterbox", pkg)
		assert.Equal(t, "github.com/tj/letterbox", mod)
		assert.Equal(t, "master", version)
		assert.Equal(t, "letterbox", bin)
	})
}

// Test parsing package paths with nested command.
func TestParsePackage_nested(t *testing.T) {
	t.Run("with version", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("https://github.com/tj/staticgen/cmd/staticgen@v1.2.0")
		assert.Equal(t, "github.com/tj/staticgen/cmd/staticgen", pkg)
		assert.Equal(t, "github.com/tj/staticgen", mod)
		assert.Equal(t, "v1.2.0", version)
		assert.Equal(t, "staticgen", bin)
	})

	t.Run("with https://github.com", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("https://github.com/tj/staticgen/cmd/staticgen")
		assert.Equal(t, "github.com/tj/staticgen/cmd/staticgen", pkg)
		assert.Equal(t, "github.com/tj/staticgen", mod)
		assert.Equal(t, "master", version)
		assert.Equal(t, "staticgen", bin)
	})

	t.Run("with github.com", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("github.com/tj/staticgen/cmd/staticgen")
		assert.Equal(t, "github.com/tj/staticgen/cmd/staticgen", pkg)
		assert.Equal(t, "github.com/tj/staticgen", mod)
		assert.Equal(t, "master", version)
		assert.Equal(t, "staticgen", bin)
	})

	t.Run("without host", func(t *testing.T) {
		pkg, mod, version, bin := parsePackage("tj/staticgen/cmd/staticgen")
		assert.Equal(t, "github.com/tj/staticgen/cmd/staticgen", pkg)
		assert.Equal(t, "github.com/tj/staticgen", mod)
		assert.Equal(t, "master", version)
		assert.Equal(t, "staticgen", bin)
	})
}
