package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"framework/fileservice/server/config"
	"framework/fileservice/server/utils"

	"github.com/sirupsen/logrus"
)

func Auth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	var jwt = config.Default.Jwt
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now   = time.Now()
			ip, _ = utils.GetIp(r)
			token = r.Header.Get("Authorization")
		)
		// "Bearer "
		if len(token) < 8 {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		var user, err = jwt.Verify(token[7:], ip)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			logrus.Error("jwt verify failed, err: %v", err)
			return
		}
		var (
			log = logrus.NewEntry(logrus.StandardLogger()).WithFields(
				logrus.Fields{
					"method": r.Method,
					"ip":     ip,
					"url":    r.RequestURI,
					"user":   user,
				})

			ctx = context.WithValue(r.Context(), "log", log)
		)
		log.Info("client request")

		handlerFunc(w, r.WithContext(ctx))

		log.WithField("latency", fmt.Sprintf("%v", time.Since(now))).
			Info("done")
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
			now   = time.Now()
			ip, _ = utils.GetIp(r)
			log   = logrus.NewEntry(logrus.StandardLogger()).WithFields(
				logrus.Fields{
					"method": r.Method,
					"ip":     ip,
					"url":    r.RequestURI,
				})

			req    request
			isJson bool
		)

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
					log.Errorf("parse form failed, err: %v", err)
					return
				}
				req.Username = r.PostForm.Get("username")
				req.Password = r.PostForm.Get("password")
			} else {
				// json提交
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "params invalid", http.StatusBadRequest)
					log.Errorf("decode login request failed, err: %v", err)
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
			token, err := jwt.CreateToken(req.Username, ip)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				log.Errorf("create token failed, err: %v", err)
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
