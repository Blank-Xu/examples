package file

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const requestTimeout = time.Second * 30

type infoResponse struct {
	Name              string `json:"name"`
	Size              int64  `json:"size"`
	ModTime           int64  `json:"mod_time"`
	UploadMaxSize     int64  `json:"upload_max_size"`
	UploadChunkSize   int64  `json:"upload_chunk_size"`
	DownloadChunkSize int64  `json:"download_chunk_size"`
}

func Info(host, filename string) (*infoResponse, error) {
	var httpClient = http.Client{Timeout: requestTimeout}
	resp, err := httpClient.Get(fmt.Sprintf("%s/info?filename=%s", host, filename))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info = new(infoResponse)
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
