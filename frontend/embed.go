// Package frontend is the frontend of the application
package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:generate pnpm i
//go:generate pnpm run build
//go:embed all:dist
var files embed.FS

func ReactHandler(path string) http.Handler {
	fsys, err := fs.Sub(files, "dist")
	if err != nil {
		log.Fatal(err)
	}

	filesystem := http.FS(fsys)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fPath := strings.TrimPrefix(r.URL.Path, path)
		checkPath := strings.TrimPrefix(fPath, "/")

		_, err := fsys.Open(checkPath)
		if os.IsNotExist(err) {
			r.URL.Path = "/"
		} else {
			r.URL.Path = fPath
		}
		http.FileServer(filesystem).ServeHTTP(w, r)
	})
}
