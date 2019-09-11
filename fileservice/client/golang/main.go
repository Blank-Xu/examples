package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"fileservice/client/golang/file"
)

var (
	action   = flag.String("a", "d", "action:[d u s i l]")
	host     = flag.String("h", "http://127.0.0.1:8080", "file server host")
	filename = flag.String("f", "1.wmv", "filename")
	count    = flag.Int("c", 1, "benchmark count")

	username = flag.String("u", "test", "username")
	password = flag.String("p", "test", "password")
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
				err = file.Download(*host, *filename, *username, *password)
			case "u": // 上传
				err = file.Upload(*host, *filename, *username, *password)
			case "s": // 删除
				err = file.Delete(*host, *filename, *username, *password)
			case "i": // 获取文件信息，包含md5值
				_, err = file.Info(*host, *filename, *username, *password, true)
			case "l": // 登录
				_, err = file.Login(*host, *username, *password)
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

	log.Printf("totalCount: %d, successCount: %d, total time cost: %v, every func cost: %v",
		*count, successCount, time.Since(now), time.Since(now).Nanoseconds()/int64(*count))
}
