package file

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func Download(host, filename string) error {
	var lfilename = filepath.Join(workDir, strconv.FormatInt(time.Now().UnixNano(), 10))

	file, err := os.OpenFile(lfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	sSize, err := InfoHead(host, filename)
	if err != nil {
		return err
	}

	var downSize int64
	for {
		info, _ := file.Stat()

		var size = sSize - info.Size()
		if size == 0 {
			break
		} else if size < 0 {
			return fmt.Errorf("local file size: %d, server file size: %d", info.Size(), sSize)
		} else if size > chunkSize {
			size = chunkSize
		}

		if size, err = downloadChunk(host, filename, file, info.Size(), info.Size()+size-1); err != nil {
			return err
		}

		downSize += size
	}

	log.Printf("download filename[%s] success, size:%d, local filename: %s\n", filename, downSize, lfilename)
	return nil
}

func downloadChunk(host, filename string, file *os.File, start, end int64) (int64, error) {
	var req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/download?filename=%s", host, filename), nil)
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

	switch resp.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
		if _, err = file.WriteAt(buf.Bytes(), start); err != nil {
			return 0, err
		}
		if err = file.Sync(); err != nil {
			return 0, err
		}
	default:
		err = errors.New(buf.String())
	}

	return end - start, err
}
