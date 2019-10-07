package http

import (
	"net/http"
)

type downloadRequest struct {
	Url string `form:"url"`
}

func download(w http.ResponseWriter, r *http.Request) {

}
