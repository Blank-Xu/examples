package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"framework/fileservice/server/config"
)

func Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			log = newLogEntry(r)
		)
		log.Info("client request")

		switch r.Method {
		case http.MethodPost, http.MethodDelete:
			var filename = r.FormValue("filename")
			if len(filename) == 0 {
				http.Error(w, "", http.StatusBadGateway)
				return
			}

			log.Infof("delete request filename: %s", filename)

			filename = filepath.Join(config.Default.FileConfig.WorkDir, filename)
			if err := os.Remove(filename); err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			log.Infof("delete file success, filename: %s", filename)

			log.WithField("latency", fmt.Sprintf("%v", time.Since(now))).
				Info("done")
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}
