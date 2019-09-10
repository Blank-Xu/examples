package file

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func Login(host, username, password string) (string, error) {
	var (
		httpClient = http.Client{Timeout: requestTimeout}
		url        = fmt.Sprintf("%s?username=%s&password=%s", host, username, password)
	)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return buf.String(), nil
	default: // 其他错误
		return "", errors.New(buf.String())
	}
}
