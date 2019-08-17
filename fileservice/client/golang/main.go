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
	action   = flag.String("a", "d", "action:[d u s i]")
	host     = flag.String("h", "http://127.0.0.1:8080", "file server host")
	filename = flag.String("f", "1.wmv", "filename")
	count    = flag.Int("c", 1, "benchmark count")
)

func main() {
	flag.Parse()

	var (
		wg  sync.WaitGroup
		now = time.Now()

		successCount int
		_lock        sync.Mutex
	)

	wg.Add(*count)

	for i := 0; i < *count; i++ {
		go func() {
			var err error
			switch *action {
			case "d": // 下载
				err = file.Download(*host, *filename)
			case "u": // 上传
				err = file.Upload(*host, *filename)
			case "s": // 删除
				err = file.Delete(*host, *filename)
			case "i": // 获取文件信息，包含md5值
				_, err = file.Info(*host, *filename, true)
			default:
				err = fmt.Errorf("not support action: %s", *action)
			}

			if err != nil {
				log.Println(err)
			} else {
				_lock.Lock()
				successCount += 1
				_lock.Unlock()
			}

			wg.Done()
		}()
	}

	wg.Wait()

	log.Printf("total time cost: %v, count: %d, every func cost: %v",
		time.Since(now), *count, time.Since(now).Nanoseconds()/int64(*count))
}
