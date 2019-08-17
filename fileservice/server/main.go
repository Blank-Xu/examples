package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"framework/fileservice/server/config"
	"framework/fileservice/server/file"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("process crashed with error, err: %v", err)
			// 等待日志记录完成
			time.Sleep(time.Second * 3)
			panic(err)
		}
	}()

	// 解析配置
	config.Init()

	// 注册路由
	http.HandleFunc("/auth", file.Auth())
	http.HandleFunc("/info", file.Info())
	http.HandleFunc("/upload", file.Upload())
	http.HandleFunc("/download", file.Download())
	http.HandleFunc("/delete", file.Delete())

	var (
		cfg = config.Default.Server

		addr = fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)

		server = http.Server{
			Addr:         addr,
			Handler:      http.DefaultServeMux,
			ReadTimeout:  time.Second * time.Duration(cfg.ReadTimeout),
			WriteTimeout: time.Second * time.Duration(cfg.WriteTimeout),
			IdleTimeout:  time.Second * time.Duration(cfg.IdleTimeout),
		}
	)

	logrus.Infof("server version: [%s], start addr: [%s]", config.VERSION, addr)

	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = cfg.MaxConnPerHost
	http.DefaultTransport.(*http.Transport).MaxIdleConns = cfg.MaxIdleConn
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = cfg.MaxIdleConnPerHost

	logrus.Warnf("server exit, %v", server.ListenAndServe())

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)

	logrus.Warnf("receive signal: %v", <-quit)
	logrus.Warn("start shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("shutdown failed: %v", err)
	}

	<-ctx.Done()
	logrus.Warnf("server exited")
}
