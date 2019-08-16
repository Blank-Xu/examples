package file

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"framework/fileservice/server/config"
)

func Upload() http.HandlerFunc {
	var uploadLimit = make(chan struct{}, config.Default.DownloadLimit)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			// auth

			var filename = r.FormValue("filename")
			if len(filename) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("params invalid"))
				return
			}

			log.Printf("upload filename: %s\n", filename)

			filename = filepath.Join(config.Default.WorkDir, filename)

			var (
				start int64
				end   int64

				err error
			)
			if _, err = fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				log.Printf("upload failed, err: %v\n", err)
				return
			}
			if end == 0 || start >= end {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("range param invalid"))
				log.Println("upload failed, err: range param invalid")
				return
			}

			var contentLength, _ = strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)

			log.Printf("upload length: %d, start: %d, end: %d\n", contentLength, start, end)

			if contentLength != (end-start+1) || contentLength > config.Default.UploadChunkSize {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("range size invalid"))
				return
			}

			var file *os.File
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				log.Printf("upload failed, open err: %v\n", err)
				return
			}
			defer file.Close()

			info, _ := file.Stat()
			if info.Size() > config.Default.UploadMaxSize {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("upload size to large, allow size: %d", config.Default.UploadMaxSize)))
				log.Printf("upload size to large, filename: %s, size: %d", info.Name(), info.Size())
				return
			}
			if info.Size() != start {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("range start invalid"))
				log.Println("upload failed, err: range start invalid")
				return
			}

			uploadLimit <- struct{}{}
			defer func() {
				<-uploadLimit
			}()

			if _, err = file.Seek(start, 2); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				log.Printf("upload failed, seek err: %v\n", err)
				return
			}

			size, err := io.CopyN(file, r.Body, contentLength)
			if err != nil && err != io.EOF {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				log.Printf("upload failed, copy err: %v\n", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strconv.FormatInt(size, 10)))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
