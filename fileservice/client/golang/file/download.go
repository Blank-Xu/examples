package file

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Download(host, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	sinfo, err := Info(host, filename)
	if err != nil {
		return err
	}

	for {
		info, err := file.Stat()
		if err != nil {
			return err
		}

		var end = sinfo.Size - info.Size()
		if end == 0 {
			log.Printf("download filename[%s] success\n", filename)
			return nil
		} else if end < 0 {
			return fmt.Errorf("local file size: %d, server file size: %d", info.Size(), sinfo.Size)
		} else if end > chunkSize {
			end = chunkSize
		}

		if err = downloadChunk(host, filename, file, info.Size(), info.Size()+end-1); err != nil {
			return err
		}
	}

	return nil
}

func downloadChunk(host, filename string, file *os.File, start, end int64) error {
	var req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/download?filename=%s", host, filename), nil)
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
	case http.StatusOK, http.StatusPartialContent:
		if _, err = file.WriteAt(buf.Bytes(), start); err != nil {
			return err
		}
		if err = file.Sync(); err != nil {
			return err
		}
		log.Println("download success, size:", end-start)
	default:
		err = errors.New(buf.String())
	}

	return err
}
