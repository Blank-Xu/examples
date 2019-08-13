package main

import (
	"fmt"
	"framework/fileservice/server/config"
	"framework/fileservice/server/file"
	"log"
	"net/http"
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

	var addr = fmt.Sprintf("%s:%d", config.Default.IP, config.Default.Port)
	log.Printf("server version: [%s], start addr: [%s]\n", config.VERSION, addr)

	log.Printf("server exit, %v", http.ListenAndServe(addr, nil))
}
