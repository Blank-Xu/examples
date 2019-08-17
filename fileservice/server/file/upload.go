package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"framework/fileservice/server/config"
)

func Upload() http.HandlerFunc {
	var (
		cfg         = config.Default.FileConfig
		uploadLimit = make(chan struct{}, cfg.DownloadLimit)
	)
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			log = newLogEntry(r)
		)
		log.Info("client request")

		switch r.Method {
		case http.MethodPost, http.MethodPut:
			// auth

			var filename = r.FormValue("filename")
			if len(filename) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("params invalid"))
				return
			}

			log.Infof("upload filename: %s", filename)

			filename = filepath.Join(cfg.WorkDir, filename)

			var (
				start int64
				end   int64

				err error
			)
			if _, err = fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				log.Errorf("upload failed, err: %v", err)
				return
			}
			if end == 0 || start >= end {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("range param invalid"))
				log.Error("upload failed, err: range param invalid")
				return
			}

			var contentLength, _ = strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)

			log.Infof("upload length: %d, start: %d, end: %d", contentLength, start, end)

			if contentLength != (end-start+1) || contentLength > cfg.UploadChunkSize {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("range size invalid"))
				log.Error("upload failed, err: range size invalid")
				return
			}

			var file *os.File
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				log.Errorf("upload failed, open filename[%s] err: %v", filename, err)
				return
			}
			defer file.Close()

			info, _ := file.Stat()
			if info.Size() > cfg.UploadMaxSize {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("upload size to large, allow size: %d", cfg.UploadMaxSize)))
				log.Errorf("upload size to large, filename: %s, size: %d", info.Name(), info.Size())
				return
			}
			if info.Size() != start {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("range start invalid"))
				log.Error("upload failed, err: range start invalid")
				return
			}

			uploadLimit <- struct{}{}
			defer func() {
				<-uploadLimit
			}()

			if _, err = file.Seek(start, 2); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				log.Errorf("upload failed, seek err: %v", err)
				return
			}

			size, err := io.CopyN(file, r.Body, contentLength)
			if err != nil && err != io.EOF {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				log.Errorf("upload failed, copy err: %v", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strconv.FormatInt(size, 10)))

			log.WithField("latency", fmt.Sprintf("%v", time.Since(now))).
				Info("done")
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
