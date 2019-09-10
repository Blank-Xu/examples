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
	requestTimeout = time.Second * 60
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

func Info(host, filename, username, password string, checkMd5 ...bool) (*InfoResponse, error) {
	var token, err = Login(host, username, password)
	if err != nil {
		return nil, err
	}

	var (
		httpClient = http.Client{Timeout: requestTimeout}
		url        string
	)
	if len(checkMd5) > 0 && checkMd5[0] == true {
		url = fmt.Sprintf("%s/info?filename=%s&md5=true", host, filename)
	} else {
		url = fmt.Sprintf("%s/info?filename=%s", host, filename)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info = new(InfoResponse)
	switch resp.StatusCode {
	case http.StatusOK:
		err = json.NewDecoder(resp.Body).Decode(info)

	case http.StatusBadGateway: // 不存在
		err = fmt.Errorf("can't find file: %s", filename)

	default: // 其他错误
		var buf bytes.Buffer
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return nil, err
		}
		err = errors.New(buf.String())
	}

	return info, err
}

func InfoHead(host, filename, token string) (int, int64, error) {
	var (
		httpClient = http.Client{Timeout: requestTimeout}
		url        = fmt.Sprintf("%s/info?filename=%s", host, filename)
	)

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return http.StatusOK, resp.ContentLength, nil
	case http.StatusNotFound:
		return http.StatusNotFound, 0, nil
	default:
		return resp.StatusCode, 0, nil
	}
}
