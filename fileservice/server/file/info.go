package file

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"framework/fileservice/server/config"
	"framework/fileservice/server/utils"
)

func Info() http.HandlerFunc {
	var (
		cfg      = config.Default.FileConfig
		md5Limit = make(chan struct{}, cfg.FileMd5Limit)
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			log = newLogEntry(r)
		)
		log.Info("client request")

		var filename = r.FormValue("filename")
		if len(filename) == 0 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		switch r.Method {
		case http.MethodHead:
			http.ServeFile(w, r, filepath.Join(cfg.WorkDir, filename))

		case http.MethodGet:
			log.Infof("info request filename: %s", filename)

			var lfilename = filepath.Join(cfg.WorkDir, filename)
			file, err := os.Stat(lfilename)
			if err != nil {
				if os.IsNotExist(err) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(http.StatusText(http.StatusNotFound)))
				} else {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
				}
				return
			}

			var md5 string
			if check := r.FormValue("md5"); check == "true" {
				md5Limit <- struct{}{}
				defer func() {
					<-md5Limit
				}()

				mfile, _ := os.OpenFile(lfilename, os.O_RDONLY, 0666)
				if mfile != nil {
					defer mfile.Close()

					md5 = utils.Md5File(mfile)
				}
			}

			var resp = infoResponse{
				Name:            file.Name(),
				Size:            file.Size(),
				ModTime:         file.ModTime().Unix(),
				Md5:             md5,
				UploadMaxSize:   cfg.UploadMaxSize,
				UploadChunkSize: cfg.UploadChunkSize,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			if err = json.NewEncoder(w).Encode(resp); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))

				log.Errorf("info json encode failed, err: %v", err)
				return
			}

			log.WithField("latency", fmt.Sprintf("%v", time.Since(now))).
				Info("done")
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

type infoResponse struct {
	Name            string `json:"name"`
	Size            int64  `json:"size"`
	ModTime         int64  `json:"mod_time"`
	Md5             string `json:"md5"`
	UploadMaxSize   int64  `json:"upload_max_size"`
	UploadChunkSize int64  `json:"upload_chunk_size"`
}
