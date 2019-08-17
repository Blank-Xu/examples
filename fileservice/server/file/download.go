package file

import (
	"fmt"
	"framework/fileservice/server/config"
	"net/http"
	"path/filepath"
	"time"
)

func Download() http.HandlerFunc {
	var downloadLimit = make(chan struct{}, config.Default.FileConfig.DownloadLimit)
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			log = newLogEntry(r)
		)
		log.Info("client request")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var filename = r.FormValue("filename")
		if len(filename) == 0 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Infof("download request filename: %s", filename)

		// check auth

		downloadLimit <- struct{}{}
		defer func() {
			<-downloadLimit
		}()

		http.ServeFile(w, r, filepath.Join(config.Default.FileConfig.WorkDir, filename))

		log.WithField("latency", fmt.Sprintf("%v", time.Since(now))).
			Info("done")
	}
}
