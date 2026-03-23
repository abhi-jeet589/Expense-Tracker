package webassets

import "embed"

// FS contains the embedded templates and static assets.
//
//go:embed templates/*.tmpl templates/transactions/*.tmpl static/*
var FS embed.FS
