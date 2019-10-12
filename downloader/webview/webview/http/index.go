package http

import (
	"net/http"

	"webview/utils"
)

const tplIndex = "views/index.html"

func index(w http.ResponseWriter, r *http.Request) {
	tpl, err := utils.OpenTpl(tplIndex)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if err = tpl.Execute(w, nil); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}
