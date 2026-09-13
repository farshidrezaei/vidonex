// Package ui provides embedded assets and utilities for the Web Studio user interface.
package ui

import "embed"

// Assets embeds the statically generated Nuxt 4 distribution for production desktop and web embedding.
//
//go:embed all:.output/public
var Assets embed.FS
