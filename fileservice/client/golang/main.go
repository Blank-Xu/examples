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
	host     = flag.String("h", "http://10.10.1.7:8080", "file server host")
	filename = flag.String("f", "tool.jar", "test filename")
)

func main() {
	flag.Parse()

	var (
		now = time.Now()
		err error
	)
	switch *action {
	case "d":
		err = file.Download(*host, *filename)
	case "u":
		err = file.Upload(*host, *filename)
	default:
		err = fmt.Errorf("not support action: %s", *action)
	}

	if err == nil {
		log.Printf("time cost: %v", time.Since(now))
	} else {
		log.Println(err)
	}
}
