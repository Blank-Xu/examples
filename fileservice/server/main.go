package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"framework/fileservice/server/config"
	"framework/fileservice/server/file"
)

func init() {
	log.SetFlags(255)
}

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

		pid = os.Getpid()

		msg = fmt.Sprintf("server pid[%d] start, version: [%s], addr: [%s]", pid, config.VERSION, addr)
	)
	logrus.Info(msg)
	log.Printf(msg)

	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = cfg.MaxConnsPerHost
	http.DefaultTransport.(*http.Transport).MaxIdleConns = cfg.MaxIdleConns
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost

	go func() {
		if err := server.ListenAndServe(); err != nil {
			msg = fmt.Sprintf("server pid[%d] exit with err: %v", pid, err)
			logrus.Error(msg)
			log.Print(msg)
		}
	}()

	var quitSignal = make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	var signalMsg = <-quitSignal
	close(quitSignal)

	msg = fmt.Sprintf("server pid[%d] received shutdown signal: [%v]", pid, signalMsg)
	logrus.Warn(msg)
	log.Print(msg)

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		msg = fmt.Sprintf("server pid[%d] shutdown failed, err: %v", pid, err)
		logrus.Error(msg)
		log.Print(msg)
	}

	<-ctx.Done()

	msg = fmt.Sprintf("server pid[%d] stoped", pid)
	logrus.Warn(msg)
	log.Print(msg)
}
