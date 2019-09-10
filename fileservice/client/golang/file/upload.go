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
)

func Upload(host, filename, username, password string, safety ...bool) error {
	file, err := os.OpenFile(filepath.Join(workDir, filename), os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	var safeUpload bool
	if len(safety) > 0 && safety[0] == true {
		safeUpload = true
	}

	token, err := Login(host, username, password)
	if err != nil {
		return err
	}

	var (
		startSize int64

		upfilename = getRandom()
	)
	if !safeUpload {
		if _, startSize, err = InfoHead(host, upfilename, token); err != nil {
			return err
		}
	}

	var (
		info, _   = file.Stat()
		uploadUrl = fmt.Sprintf("%s/upload?filename=%s", host, upfilename)

		upSize int64
	)

	for {
		if safeUpload {
			if _, startSize, err = InfoHead(host, upfilename, token); err != nil {
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

		if size, err = uploadChunk(uploadUrl, token, bytes.NewReader(data), startSize, startSize-1+int64(size)); err != nil {
			return err
		}

		startSize += size
		upSize += size
	}

	log.Printf("upload file[%s] success, upfilename: %s size: %d", filename, upfilename, upSize)
	return nil
}

func uploadChunk(url, token string, body io.Reader, start, end int64) (int64, error) {
	var req, _ = http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	req.Header.Set("Authorization", "Bearer "+token)

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
