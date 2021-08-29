package images

import (
	"embed"
	"net/http"
)

//go:embed *.svg
var fs embed.FS

func NewHandler() http.Handler {
	return http.FileServer(http.FS(fs))
}
