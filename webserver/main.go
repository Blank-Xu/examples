package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"webserver/init"
	"webserver/tools"
	"webserver/tools/httpserver"
)

func main() {
	pid := os.Getpid()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("server pid[%d] crashed with err: %v", pid, err)
			time.Sleep(time.Second)
			panic(err)
		}
		tools.Logger.Sync()
	}()

	// 解析配置
	init.Init()

	// 注册路由

	server := httpserver.Default.NewHttpServer(http.DefaultServeMux)

	log.Printf("server pid[%d] start, addr: [%s]", pid, server.Addr)
	go func() {

		if err := server.ListenAndServe(); err != nil {
			log.Printf("server pid[%d] exit with err: %v", pid, err)
		}
	}()

	httpserver.ShutDown(server)
}
