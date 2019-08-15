package main

import (
	"context"
	"fmt"
	"framework/fileservice/server/config"
	"framework/fileservice/server/file"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("process exit with error, err: %v\n", err)
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

	var (
		addr = fmt.Sprintf("%s:%d", config.Default.IP, config.Default.Port)

		server = http.Server{
			Addr:         addr,
			Handler:      http.DefaultServeMux,
			ReadTimeout:  time.Second * 30,
			WriteTimeout: time.Second * 30,
			IdleTimeout:  time.Second * 90,
		}
	)

	log.Printf("server version: [%s], start addr: [%s]\n", config.VERSION, addr)

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	log.Printf("server exit, %v", server.ListenAndServe())

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)
	log.Printf("recive signal: %v\n", <-quit)

	log.Println("start shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}

	<-ctx.Done()
	log.Println("server exited")
}
