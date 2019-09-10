package file

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
)

func getRandom() (str string) {
	for i := 0; i < 15; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(100))
		str += fmt.Sprintf("%d", n)
	}
	return
}

func Download(host, filename, username, password string) error {
	var lfilename = filepath.Join(workDir, getRandom())

	file, err := os.OpenFile(lfilename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	token, err := Login(host, username, password)
	if err != nil {
		return err
	}

	statusCode, sSize, err := InfoHead(host, filename, token)
	if err != nil {
		return err
	}
	if statusCode == http.StatusNotFound {
		return errors.New("server file not found")
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

		if size, err = downloadChunk(host, filename, token, file, info.Size(), info.Size()+size-1); err != nil {
			return err
		}

		downSize += size
	}

	log.Printf("download filename[%s] success, size:%d, local filename: %s", filename, downSize, lfilename)
	return nil
}

func downloadChunk(host, filename, token string, file *os.File, start, end int64) (int64, error) {
	var req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/download?filename=%s", host, filename), nil)
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
