package file

import (
	"bytes"
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
			Name:    file.Name(),
			Size:    file.Size(),
			ModTime: file.ModTime().Unix(),
		}

		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			log.Printf("info json encode failed, err: %v\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(buf.Bytes()); err != nil {
			log.Printf("info response data failed, err: %v\n", err)
		}
	}
}

type infoResponse struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
}
