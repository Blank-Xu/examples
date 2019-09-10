package file

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func Delete(host, filename, username, password string) error {
	token, err := Login(host, username, password)
	if err != nil {
		return err
	}

	var (
		httpClient = http.Client{Timeout: requestTimeout}
		url        = fmt.Sprintf("%s/delete?filename=%s", host, filename)
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default: // 其他错误
		var buf bytes.Buffer
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return err
		}
		return errors.New(buf.String())
	}
}
