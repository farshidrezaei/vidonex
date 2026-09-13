// Package ui provides embedded assets and utilities for the Web Studio user interface.
package ui

import (
	"embed"
	"io/fs"
)

// EmbeddedFS embeds the static web assets from the public directory.
//
//go:embed all:public
var EmbeddedFS embed.FS

// GetFS returns an fs.FS pointing to the static web application root.
func GetFS() (fs.FS, error) {
	// If generated dist exists under public, serve it; otherwise serve public root
	if sub, err := fs.Sub(EmbeddedFS, "public/dist"); err == nil {
		if _, statErr := fs.Stat(sub, "index.html"); statErr == nil {
			return sub, nil
		}
	}
	return fs.Sub(EmbeddedFS, "public")
}
