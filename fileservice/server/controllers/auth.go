package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"framework/fileservice/server/config"

	"github.com/sirupsen/logrus"
)

func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		http.Error(w, "", http.StatusUnauthorized)
	}
}

func Login() http.HandlerFunc {
	type request struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	var jwt = config.Default.Jwt

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			log = newLogEntry(r)

			req    request
			isJson bool
		)
		log.Info("login request")

		switch r.Method {
		case http.MethodGet:
			// url提交
			req.Username = r.FormValue("username")
			req.Password = r.FormValue("password")
		case http.MethodPost:
			// form提交
			if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				if err := r.ParseForm(); err != nil {
					http.Error(w, "params invalid", http.StatusBadRequest)
					logrus.Errorf("parse form failed, err: %v", err)
					return
				}
				req.Username = r.PostForm.Get("username")
				req.Password = r.PostForm.Get("password")
			} else {
				// json提交
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "params invalid", http.StatusBadRequest)
					logrus.Errorf("decode login request failed, err: %v", err)
					return
				}
				isJson = true
			}
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		if len(req.Username) == 0 || len(req.Password) == 0 {
			http.Error(w, "params invalid", http.StatusBadRequest)
			return
		}

		if req.Username == "test" && req.Password == "test" {
			token, err := jwt.CreateToken(req.Username)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				logrus.Errorf("create token failed, err: %v", err)
				return
			}

			if isJson {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"token":"%s"}`, token)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(token))
		}

		log.WithField("latency", time.Since(now)).Info("done")
	}
}
