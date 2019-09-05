package file

import (
	"net/http"
	"time"
)

func Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			log = newLogEntry(r)
		)
		log.Info("client request")

		switch r.Method {
		case http.MethodGet:

		case http.MethodPost:

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		log.WithField("latency", time.Since(now)).Info("done")
	}
}
