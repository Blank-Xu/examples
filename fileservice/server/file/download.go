package file

import (
	"framework/fileservice/server/config"
	"log"
	"net/http"
	"path/filepath"
)

func Download() http.HandlerFunc {
	var downloadLimit = make(chan struct{}, config.Default.DownloadLimit)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var filename = r.FormValue("filename")
		if len(filename) == 0 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Printf("download request filename: %s\n", filename)

		// check auth

		downloadLimit <- struct{}{}
		defer func() {
			<-downloadLimit
		}()

		http.ServeFile(w, r, filepath.Join(config.Default.WorkDir, filename))
	}
}
