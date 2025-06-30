package web

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
)

//go:embed static/*.js static/*.css static/*.ico static/*.png static/site.webmanifest
var WebStaticFS embed.FS

func AssetPath(fsys fs.FS, path string) string {
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return "/static/" + path
	}
	sum := sha256.Sum256(b)
	return fmt.Sprintf("/static/%s?v=%x", path, sum[:8])
}
