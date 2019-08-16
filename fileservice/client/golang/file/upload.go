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

func Upload(host, filename string, safety ...bool) error {
	file, err := os.OpenFile(filepath.Join(workDir, filename), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	var safeUpload bool
	if len(safety) > 0 && safety[0] == true {
		safeUpload = true
	}

	var (
		startSize  int64
		upfilename = strconv.FormatInt(time.Now().UnixNano(), 10)
	)
	if !safeUpload {
		if startSize, err = infoSize(host, upfilename); err != nil {
			return err
		}
	}

	var (
		info, _ = file.Stat()
		upSize  int64
	)

	for {
		if safeUpload {
			if startSize, err = infoSize(host, upfilename); err != nil {
				return err
			}
		}

		var size = info.Size() - startSize
		if size == 0 { // 已上传完成
			break
		} else if size < 0 { // 服务器文件大小错误
			return fmt.Errorf("server file size error, filename: %s local size: %d, server size: %d", upfilename, info.Size(), startSize)
		} else if size > chunkSize { // 检查是否需要分包上传
			size = chunkSize
		}

		var data = make([]byte, size)
		if _, err = file.ReadAt(data, startSize); err != nil && err != io.EOF {
			return err
		}

		if size, err = uploadChunk(host, upfilename, bytes.NewReader(data), startSize, startSize-1+int64(size)); err != nil {
			return err
		}

		startSize += size
		upSize += size
	}

	log.Printf("upload success, size: %d", upSize)
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

func infoSize(host, upfilename string) (int64, error) {
	sinfo, err := Info(host, upfilename)
	if err != nil && err != os.ErrExist {
		return 0, err
	}
	if sinfo == nil {
		return 0, errors.New("info data nil")
	}
	return sinfo.Size, nil
}
