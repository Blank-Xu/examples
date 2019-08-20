package file

import (
	"framework/fileservice/server/utils"
	"net/http"

	"github.com/sirupsen/logrus"
)

func newLogEntry(r *http.Request) *logrus.Entry {
	var ip, _ = utils.GetIp(r)

	var fields = logrus.Fields{
		"method":    r.Method,
		"ip":        ip,
		"url":       r.RequestURI,
		"user_id":   "",
		"user_name": "",
	}

	return logrus.NewEntry(logrus.StandardLogger()).WithFields(fields)
}
