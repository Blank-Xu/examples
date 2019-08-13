package utils

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrCheckSignFailed = errors.New("check sign failed")

func NewUpload(r *http.Request, chunkSize int64) *Upload {
	return &Upload{Request: r, ChunkSize: chunkSize}
}

type Upload struct {
	*http.Request
	ChunkSize int64

	FileName string

	Length int64
	Start  int64
	End    int64
}

func (p *Upload) Parse() (err error) {
	_, err = fmt.Sscanf(p.Header.Get("Content-Length"), "%d", &p.Length)
	_, err = fmt.Sscanf(p.Header.Get("Content-Upload"), "range=%d-%d", &p.Start, &p.End)
	return err
}

func (p *Upload) checkSign() bool {

	return true
}
