package file

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	requestTimeout = time.Second * 30
	chunkSize      = 4 * 1024 * 1024
	workDir        = "files"
)

func init() {
	os.MkdirAll(workDir, 0666)
}

type InfoResponse struct {
	Name              string `json:"name"`
	Size              int64  `json:"size"`
	ModTime           int64  `json:"mod_time"`
	Md5               string `json:"md5"`
	UploadMaxSize     int64  `json:"upload_max_size"`
	UploadChunkSize   int64  `json:"upload_chunk_size"`
	DownloadChunkSize int64  `json:"download_chunk_size"`
}

func Info(host, filename string, checkMd5 ...bool) (*InfoResponse, error) {
	var (
		httpClient = http.Client{Timeout: requestTimeout}
		url        string
	)
	if len(checkMd5) > 0 && checkMd5[0] == true {
		url = fmt.Sprintf("%s/info?filename=%s&md5=true", host, filename)
	} else {
		url = fmt.Sprintf("%s/info?filename=%s", host, filename)
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info = new(InfoResponse)
	switch resp.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(resp.Body).Decode(info)

	case http.StatusBadGateway: // 不存在

	default: // 其他错误
		var buf bytes.Buffer
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return nil, err
		}
		err = errors.New(buf.String())
	}

	return info, err
}

func InfoHead(host, filename string) (int64, error) {
	var (
		httpClient = http.Client{Timeout: requestTimeout}
		url        = fmt.Sprintf("%s/info?filename=%s", host, filename)
	)
	resp, err := httpClient.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.ContentLength, nil
	default:
		return 0, errors.New("404")
	}
}