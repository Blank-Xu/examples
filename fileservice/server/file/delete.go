package file

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
				w.WriteHeader(http.StatusBadGateway)
				return
			}

			log.Infof("delete request filename: %s", filename)

			// check auth

			filename = filepath.Join(config.Default.FileConfig.WorkDir, filename)
			if err := os.Remove(filename); err != nil {
				if os.IsNotExist(err) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(http.StatusText(http.StatusNotFound)))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				}
				return
			}
			log.Infof("delete file success, filename: %s", filename)

			log.WithField("latency", fmt.Sprintf("%v", time.Since(now))).
				Info("done")
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
