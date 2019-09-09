package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

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

			var log = r.Context().Value("log").(*logrus.Entry)
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
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}
