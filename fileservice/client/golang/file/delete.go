package file

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func Delete(host, filename string) error {
	var (
		httpClient = http.Client{Timeout: requestTimeout}

		url = fmt.Sprintf("%s/delete?filename=%s", host, filename)
	)
	resp, err := httpClient.PostForm(url, nil)
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
