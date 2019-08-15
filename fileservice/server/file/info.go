package file

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"framework/fileservice/server/config"
)

func Info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		var filename = r.FormValue("filename")
		if len(filename) == 0 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Printf("info request filename: %s\n", filename)

		file, err := os.Stat(filepath.Join(config.Default.WorkDir, filename))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var resp = infoResponse{
			Name:              file.Name(),
			Size:              file.Size(),
			ModTime:           file.ModTime().Unix(),
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
	UploadMaxSize     int64  `json:"upload_max_size"`
	UploadChunkSize   int64  `json:"upload_chunk_size"`
	DownloadChunkSize int64  `json:"download_chunk_size"`
}
