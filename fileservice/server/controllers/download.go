package controllers

import (
	"fmt"
	"framework/fileservice/server/config"
	"net/http"
	"path/filepath"
	"time"
)

func Download() http.HandlerFunc {
	var (
		cfg           = config.Default.FileConfig
		downloadLimit = make(chan struct{}, cfg.DownloadLimit)
	)

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			log = newLogEntry(r)
		)
		log.Info("client request")

		if r.Method != http.MethodGet {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var filename = r.FormValue("filename")
		if len(filename) == 0 {
			http.Error(w, "", http.StatusBadGateway)
			return
		}

		log.Infof("download request filename: %s", filename)

		downloadLimit <- struct{}{}
		defer func() {
			<-downloadLimit
		}()

		http.ServeFile(w, r, filepath.Join(cfg.WorkDir, filename))

		log.WithField("latency", fmt.Sprintf("%v", time.Since(now))).
			Info("done")
	}
}
