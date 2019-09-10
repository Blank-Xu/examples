package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"framework/fileservice/server/config"
	"framework/fileservice/server/utils"
)

func Info() http.HandlerFunc {
	type response struct {
		Name            string `json:"name"`
		Size            int64  `json:"size"`
		ModTime         int64  `json:"mod_time"`
		Md5             string `json:"md5"`
		UploadMaxSize   int64  `json:"upload_max_size"`
		UploadChunkSize int64  `json:"upload_chunk_size"`
	}

	var (
		cfg      = config.Default.FileConfig
		md5Limit = make(chan struct{}, cfg.FileMd5Limit)
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var filename = r.FormValue("filename")
		if len(filename) == 0 {
			http.Error(w, "", http.StatusBadGateway)
			return
		}

		switch r.Method {
		case http.MethodHead:
			http.ServeFile(w, r, filepath.Join(cfg.WorkDir, filename))

		case http.MethodGet:
			var ctx = r.Context().Value(ContextKey).(*ContextValue)
			ctx.Log.Infof("info request filename: %s", filename)

			var lfilename = filepath.Join(cfg.WorkDir, filename)
			file, err := os.Stat(lfilename)
			if err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
				} else {
					http.Error(w, err.Error(), http.StatusBadRequest)
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
				if mfile == nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					ctx.Log.Error(err)
					return
				}
				defer mfile.Close()

				md5, err = utils.Md5File(mfile)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					ctx.Log.Error(err)
					return
				}
			}

			var resp = response{
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				ctx.Log.Errorf("info json encode failed, err: %v", err)
				return
			}
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}
