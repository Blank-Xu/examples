package main

import (
	"flag"
	"fmt"
	"framework/fileservice/client/golang/file"
	"log"
	"time"
)

var (
	action   = flag.String("a", "d", "upload")
	host     = flag.String("h", "http://127.0.0.1:8080", "file server host")
	filename = flag.String("f", "1.wmv", "test filename")
)

func main() {
	flag.Parse()

	var (
		now  = time.Now()
		info *file.InfoResponse
		err  error
	)
	switch *action {
	case "d":
		err = file.Download(*host, *filename)
	case "u":
		err = file.Upload(*host, *filename)
	case "i":
		info, err = file.Info(*host, *filename, true)
	default:
		err = fmt.Errorf("not support action: %s", *action)
	}

	if err == nil {
		log.Printf("time cost: %v", time.Since(now))
		if info != nil {
			log.Println(info)
		}
	} else {
		log.Println(err)
	}
}
