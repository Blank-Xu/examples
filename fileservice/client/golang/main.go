package main

import (
	"flag"
	"fmt"
	"framework/fileservice/client/golang/file"
	"log"
	"sync"
	"time"
)

var (
	action   = flag.String("a", "d", "upload")
	host     = flag.String("h", "http://127.0.0.1:8080", "file server host")
	filename = flag.String("f", "1.wmv", "filename")
	count    = flag.Int("c", 1, "benchmark count")
)

func main() {
	flag.Parse()

	var (
		wg  sync.WaitGroup
		now = time.Now()
	)

	wg.Add(*count)

	for i := 0; i < *count; i++ {
		go func() {
			var err error
			switch *action {
			case "d":
				err = file.Download(*host, *filename)
			case "u":
				err = file.Upload(*host, *filename)
			case "i":
				_, err = file.Info(*host, *filename, true)
			default:
				err = fmt.Errorf("not support action: %s", *action)
			}

			if err != nil {
				log.Println(err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	log.Printf("total time cost: %v, count: %d, every func cost: %v",
		time.Since(now), *count, float64(time.Since(now).Nanoseconds())/float64(*count))
}
