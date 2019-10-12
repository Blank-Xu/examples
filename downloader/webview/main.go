package main

import (
	"log"

	"github.com/zserge/webview"

	"webview/config"
	"webview/webview/http"
	"webview/webview/rpc"
)

const (
	windowWidth  = 480
	windowHeight = 320
)

func main() {
	url := http.Init()
	log.Println(url)

	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  config.Title,
		Resizable:              true,
		URL:                    url,
		ExternalInvokeCallback: rpc.Init,
	})
	w.SetColor(255, 255, 255, 255)
	defer w.Exit()
	w.Run()
}
