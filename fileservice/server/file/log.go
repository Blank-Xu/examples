package file

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func newLogEntry(r *http.Request) *logrus.Entry {
	var (
		addr = strings.Split(r.RemoteAddr, ":")
		ip   string
	)
	if len(addr) > 0 {
		ip = addr[0]
	}

	fields := logrus.Fields{
		"method":    r.Method,
		"ip":        ip,
		"url":       r.RequestURI,
		"user_id":   "",
		"user_name": "",
	}

	return logrus.NewEntry(logrus.StandardLogger()).WithFields(fields)
}
