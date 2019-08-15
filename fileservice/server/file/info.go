package file

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"framework/fileservice/server/config"
	"framework/fileservice/server/utils"
)

func Info() http.HandlerFunc {
	var md5Limit = make(chan struct{}, config.Default.FileMd5Limit)
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

		log.Printf("info request filename: %s\n", filename)

		var lfilename = filepath.Join(config.Default.WorkDir, filename)
		file, err := os.Stat(lfilename)
		if err != nil {
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusBadGateway)
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
			md5 = utils.Md5File(mfile)
		}

		var resp = infoResponse{
			Name:              file.Name(),
			Size:              file.Size(),
			ModTime:           file.ModTime().Unix(),
			Md5:               md5,
			UploadMaxSize:     config.Default.UploadMaxSize,
			UploadChunkSize:   config.Default.UploadChunkSize,
			DownloadChunkSize: config.Default.DownloadChunkSize,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			log.Printf("info json encode failed, err: %v\n", err)
		}
	}
}

type infoResponse struct {
	Name              string `json:"name"`
	Size              int64  `json:"size"`
	ModTime           int64  `json:"mod_time"`
	Md5               string `json:"md5"`
	UploadMaxSize     int64  `json:"upload_max_size"`
	UploadChunkSize   int64  `json:"upload_chunk_size"`
	DownloadChunkSize int64  `json:"download_chunk_size"`
}
