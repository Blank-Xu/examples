package controllers

import (
	"net/http"
	"path/filepath"

	"fileservice/server/config"
)

func Download() http.HandlerFunc {
	var (
		cfg           = config.Default.FileConfig
		downloadLimit = make(chan struct{}, cfg.DownloadLimit)
	)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var filename = r.FormValue("filename")
		if len(filename) == 0 {
			http.Error(w, "", http.StatusBadGateway)
			return
		}

		var ctx = r.Context().Value(ContextKey).(*ContextValue)
		ctx.Log.Infof("download request filename: %s", filename)

		// TODO: 检查 ctx.User 是否有下载权限

		downloadLimit <- struct{}{}
		defer func() {
			<-downloadLimit
		}()

		http.ServeFile(w, r, filepath.Join(cfg.WorkDir, filename))
	}
}
