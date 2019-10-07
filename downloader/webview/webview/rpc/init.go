package rpc

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/zserge/webview"
)

var routerMap = map[string]func(webview.WebView){
	"download": download,
	// "counter":  counter,
}

func Init(w webview.WebView, data string) {
	rpc, ok := routerMap[data]
	fmt.Println(data)
	if !ok {
		return
	}
	rpc(w)
	return

	switch {
	case data == "close":
		w.Terminate()
	case data == "fullscreen":
		w.SetFullscreen(true)
	case data == "unfullscreen":
		w.SetFullscreen(false)
	case data == "open":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, 0, "Open file", ""))
	case data == "opendir":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, webview.DialogFlagDirectory, "Open directory", ""))
	case data == "save":
		log.Println("save", w.Dialog(webview.DialogTypeSave, 0, "Save file", ""))
	case data == "message":
		w.Dialog(webview.DialogTypeAlert, 0, "Hello", "Hello, world!")
	case data == "info":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "Hello", "Hello, info!")
	case data == "warning":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagWarning, "Hello", "Hello, warning!")
	case data == "error":
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagError, "Hello", "Hello, error!")
	case strings.HasPrefix(data, "changeTitle:"):
		w.SetTitle(strings.TrimPrefix(data, "changeTitle:"))
	case strings.HasPrefix(data, "changeColor:"):
		hex := strings.TrimPrefix(strings.TrimPrefix(data, "changeColor:"), "#")
		num := len(hex) / 2
		if !(num == 3 || num == 4) {
			log.Println("Color must be RRGGBB or RRGGBBAA")
			return
		}
		i, err := strconv.ParseUint(hex, 16, 64)
		if err != nil {
			log.Println(err)
			return
		}
		if num == 3 {
			r := uint8((i >> 16) & 0xFF)
			g := uint8((i >> 8) & 0xFF)
			b := uint8(i & 0xFF)
			w.SetColor(r, g, b, 255)
			return
		}
		if num == 4 {
			r := uint8((i >> 24) & 0xFF)
			g := uint8((i >> 16) & 0xFF)
			b := uint8((i >> 8) & 0xFF)
			a := uint8(i & 0xFF)
			w.SetColor(r, g, b, a)
			return
		}
	}
}
