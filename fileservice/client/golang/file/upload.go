package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const chunkSize = 4 * 1024 * 1024

func Upload(host, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	info, _ := file.Stat()

	var upfilename = strconv.FormatInt(time.Now().UnixNano(), 10)

	for {
		sinfo, err := Info(host, upfilename)
		if err != nil {
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
		if err = uploadChunk(host, upfilename, bytes.NewReader(data), sinfo.Size-1, sinfo.Size-1+int64(size)); err != nil {
			return err
		}
	}
}

func uploadChunk(host, filename string, body io.Reader, start, end int64) error {
	var req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/upload?filename=%s", host, filename), body)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	var httpClient = http.Client{Timeout: requestTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		log.Println("upload success, size:", buf.String())
	default:
		err = errors.New(buf.String())
	}
	return err
}
