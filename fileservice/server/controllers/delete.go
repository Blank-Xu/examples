package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"framework/fileservice/server/config"
)

func Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodDelete:
			var filename = r.FormValue("filename")
			if len(filename) == 0 {
				http.Error(w, "", http.StatusBadGateway)
				return
			}

			var ctx = r.Context().Value(ContextKey).(*ContextValue)
			ctx.Log.Infof("delete request filename: %s", filename)

			// TODO: 检查 ctx.User 是否有删除权限
			filename = filepath.Join(config.Default.FileConfig.WorkDir, filename)
			if err := os.Remove(filename); err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			ctx.Log.Infof("delete file success, filename: %s", filename)
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}
