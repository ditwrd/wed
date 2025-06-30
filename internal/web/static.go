package web

import "embed"

//go:embed static/*.js static/*.css
var WebStaticFS embed.FS
