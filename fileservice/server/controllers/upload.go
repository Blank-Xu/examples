package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"framework/fileservice/server/config"
)

func Upload() http.HandlerFunc {
	var (
		cfg         = config.Default.FileConfig
		uploadLimit = make(chan struct{}, cfg.DownloadLimit)
	)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var filename = r.FormValue("filename")
			if len(filename) == 0 {
				http.Error(w, "params invalid", http.StatusBadRequest)
				return
			}

			var ctx = r.Context().Value(ContextKey).(*ContextValue)
			ctx.Log.Infof("upload filename: %s", filename)

			// TODO: 检查 ctx.User 是否有上传权限

			filename = filepath.Join(cfg.WorkDir, filename)
			var (
				start int64
				end   int64

				err error
			)
			if _, err = fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				ctx.Log.Errorf("upload failed, err: %v", err)
				return
			}
			if end == 0 || start >= end {
				http.Error(w, "range param invalid", http.StatusBadRequest)
				ctx.Log.Error("upload failed, err: range param invalid")
				return
			}

			var contentLength, _ = strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)

			ctx.Log.Infof("upload length: %d, start: %d, end: %d", contentLength, start, end)

			if contentLength != (end-start+1) || contentLength > cfg.UploadChunkSize {
				http.Error(w, "range size invalid", http.StatusBadRequest)
				ctx.Log.Error("upload failed, err: range size invalid")
				return
			}

			var file *os.File
			file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				ctx.Log.Errorf("upload failed, open filename[%s] err: %v", filename, err)
				return
			}
			defer file.Close()

			info, _ := file.Stat()
			if info.Size() > cfg.UploadMaxSize {
				http.Error(w, fmt.Sprintf("upload size to large, allow size: %d", cfg.UploadMaxSize), http.StatusBadRequest)
				ctx.Log.Errorf("upload size to large, filename: %s, size: %d", info.Name(), info.Size())
				return
			}
			if info.Size() != start {
				http.Error(w, "range start invalid", http.StatusBadRequest)
				ctx.Log.Error("upload failed, err: range start invalid")
				return
			}

			uploadLimit <- struct{}{}
			defer func() {
				<-uploadLimit
			}()

			if _, err = file.Seek(start, 2); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				ctx.Log.Errorf("upload failed, seek err: %v", err)
				return
			}

			size, err := io.CopyN(file, r.Body, contentLength)
			if err != nil && err != io.EOF {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				ctx.Log.Errorf("upload failed, copy err: %v", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strconv.FormatInt(size, 10)))
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}
