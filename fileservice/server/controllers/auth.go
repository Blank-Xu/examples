package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fileservice/server/config"
	"fileservice/server/utils"

	"github.com/sirupsen/logrus"
)

const ContextKey = "CtxKey"

type ContextValue struct {
	User string
	Ip   string
	Log  *logrus.Entry
}

// TODO: 可以加入 ip 请求次数检测，ip 白名单和黑名单验证

func Auth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	jwt := config.Default.Jwt

	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ip, _ := utils.GetIp(r)
		token := r.Header.Get("Authorization")

		// "Bearer "
		if len(token) < 8 {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		user, err := jwt.Verify(token[7:], ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			logrus.Error(err)
			return
		}

		log := logrus.NewEntry(logrus.StandardLogger()).WithFields(
			logrus.Fields{
				"method": r.Method,
				"ip":     ip,
				"url":    r.URL.Path,
				"user":   user,
			})

		ctx := context.WithValue(r.Context(), ContextKey, &ContextValue{
			User: user,
			Ip:   ip,
			Log:  log,
		})

		handlerFunc(w, r.WithContext(ctx))

		log.WithField("latency", float64(time.Now().Sub(now).Nanoseconds())/1000000.0).Info("done")
	}
}

func Login() http.HandlerFunc {
	type request struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	jwt := config.Default.Jwt

	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		ip, _ := utils.GetIp(r)
		log := logrus.NewEntry(logrus.StandardLogger()).WithFields(
			logrus.Fields{
				"method": r.Method,
				"ip":     ip,
				"url":    r.URL.Path,
			})

		var (
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

		if req.Username == "" || req.Password == "" {
			http.Error(w, "params invalid", http.StatusBadRequest)
			return
		}

		log = log.WithField("user", req.Username)

		// TODO: 用户验证需要从数据库或其他配置表重新读取
		if req.Username == jwt.Username && req.Password == jwt.Password {
			token, err := jwt.CreateToken(req.Username, ip)
			if err != nil {
				http.Error(w, "create token failed", http.StatusInternalServerError)
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

		log.WithField("latency", float64(time.Now().Sub(now).Nanoseconds())/1000000.0).Info("done")
	}
}
