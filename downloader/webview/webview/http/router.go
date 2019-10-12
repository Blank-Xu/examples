package http

import (
	"net/http"
)

func router() http.Handler {
	handler := http.NewServeMux()

	handler.Handle("/statics/",
		http.StripPrefix("/statics/",
			http.FileServer(http.Dir("statics"))))

	handler.HandleFunc("/", index)
	handler.HandleFunc("/download", download)

	return handler
}
