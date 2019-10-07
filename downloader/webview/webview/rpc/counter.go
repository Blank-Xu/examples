package rpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/zserge/webview"
)

var count = 0
var ch = make(chan int, 1)
var so = sync.Once{}

type Counter struct {
}

func (p *Counter) Add(w webview.WebView) {
	so.Do(func() {
		go doCount()
	})

	// for v := range ch {
	//
	// }
	fmt.Println("download", w)
}

func doCount() {
	count++
	ch <- count
	time.Sleep(time.Second)
}
