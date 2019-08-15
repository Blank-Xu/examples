package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Upload(host, filename string) error {
	file, err := os.OpenFile(filepath.Join(workDir, filename), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	var (
		info, _    = file.Stat()
		upfilename = strconv.FormatInt(time.Now().UnixNano(), 10)

		upsize int64
	)
	for {
		sinfo, err := Info(host, upfilename)
		if err != nil && err != os.ErrExist {
			return err
		}

		var size = info.Size() - sinfo.Size
		if size == 0 { // 已上传完成
			return nil
		} else if size < 0 { // 服务器文件大小错误
			return fmt.Errorf("server file size error, filename: %s local size: %d, server size: %d", upfilename, info.Size(), sinfo.Size)
		} else if size > chunkSize { // 检查是否需要分包上传
			size = chunkSize
		}

		var data = make([]byte, size)
		if _, err = file.ReadAt(data, sinfo.Size); err != nil && err != io.EOF {
			return err
		}
		if size, err = uploadChunk(host, upfilename, bytes.NewReader(data), sinfo.Size, sinfo.Size-1+int64(size)); err != nil {
			return err
		}
		upsize += size
	}
	log.Printf("upload success, size: %d", upsize)
	return nil
}

func uploadChunk(host, filename string, body io.Reader, start, end int64) (int64, error) {
	var req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/upload?filename=%s", host, filename), body)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	var httpClient = http.Client{Timeout: requestTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return 0, err
	}

	var size int64
	switch resp.StatusCode {
	case http.StatusOK:
		size, _ = strconv.ParseInt(buf.String(), 10, 64)
	default:
		err = errors.New(buf.String())
	}
	return size, err
}
