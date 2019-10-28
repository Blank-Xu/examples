package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"ftp"
)

var (
	configFile = flag.String("s", "config.yaml", "config file")
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	log.Println("server starting...")

	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("server crashed with error: %v", err)
			// 等待日志记录完成
			time.Sleep(time.Second)
			panic(err)
		}
		time.Sleep(time.Second)
	}()

	// 解析配置
	var cfg ftp.Config
	file, err := os.Open(*configFile)
	if err != nil {
		log.Println(err)
		log.Println("use default config")
	} else if file != nil {
		if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
			log.Fatalf("parse config file [%s] failed, err: %v\n", *configFile, err)
		}
		log.Printf("load config file [%s] success\n", *configFile)
	}

	server, err := ftp.NewServer(&cfg)
	if err != nil {
		log.Fatalf("server create failed, err: %v\n", err)
	}

	log.Println("server started")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("server exit with err: %v\n", err)
			}
		}
	}()

	var quitSignal = make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	log.Printf("server receive shutdown signal: [%v]\n", <-quitSignal)

	if err = server.Close(); err != nil {
		log.Printf("server closed with error: %v\n", err)
	}

	log.Println("server exited.")
}
