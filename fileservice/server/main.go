package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fileservice/server/config"
	"fileservice/server/controllers"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	pid := os.Getpid()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("server pid[%d] crashed with error: %v", pid, err)
			// 等待日志记录完成
			time.Sleep(time.Second)
			panic(err)
		}
		time.Sleep(time.Second)
	}()

	// 解析配置
	config.Init()

	// 注册路由
	http.HandleFunc("/", controllers.Login())
	http.HandleFunc("/info", controllers.Auth(controllers.Info()))
	http.HandleFunc("/upload", controllers.Auth(controllers.Upload()))
	http.HandleFunc("/download", controllers.Auth(controllers.Download()))
	http.HandleFunc("/delete", controllers.Auth(controllers.Delete()))

	server := config.Default.Server.NewServer(http.DefaultServeMux)

	go func() {
		log.Printf("server pid[%d] start, version: [%s], addr: [%s]", pid, config.VERSION, server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server pid[%d] exit with err: %v", pid, err)
		}
	}()

	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	log.Printf("server pid[%d] receive shutdown signal: [%v]", pid, <-quitSignal)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server pid[%d] shutdown failed, err: %v", pid, err)
	}

	log.Printf("server pid[%d] stoped", pid)
}
