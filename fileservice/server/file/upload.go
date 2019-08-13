package file

import (
	"io"
	"log"
	"os"

	"framework/fileservice/server/config"
	"framework/fileservice/server/utils"

	"net/http"
)

func Upload() http.HandlerFunc {
	var uploadLimit = make(chan struct{}, config.Default.DownloadLimit)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			uploadLimit <- struct{}{}
			defer func() {
				<-uploadLimit
			}()

			var (
				upload = utils.NewUpload(r, config.Default.UploadChunkSize)
				err    error
			)
			if err = upload.Parse(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			log.Printf("upload request filename: %s\n", upload.FileName)

			file, err := os.OpenFile(upload.FileName, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			defer file.Close()

			if _, err = file.Seek(upload.Start, 2); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, config.Default.UploadChunkSize)

			reader, err := r.MultipartReader()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			for {
				form, err := reader.ReadForm(config.Default.UploadChunkSize)
				if err == io.EOF {
					break
				}
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}

			}

		default:
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
