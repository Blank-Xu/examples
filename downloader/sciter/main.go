package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

func main() {
	// 创建window窗口
	// 参数一表示创建窗口的样式
	// SW_TITLEBAR 顶层窗口，有标题栏
	// SW_RESIZEABLE 可调整大小
	// SW_CONTROLS 有最小/最大按钮
	// SW_MAIN 应用程序主窗口，关闭后其他所有窗口也会关闭
	// SW_ENABLE_DEBUG 可以调试
	// 参数二表示创建窗口的矩形
	path, err := os.Getwd()
	// path,err := utils.GetPath()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(path)

	// path = filepath.Join(path,"sciter.dll")
	// log.Println(path)

	// sciter.SetDLL("sciter.dll")
	w, err := window.New(sciter.SW_TITLEBAR|
		sciter.SW_RESIZEABLE|
		sciter.SW_CONTROLS|
		sciter.SW_MAIN|
		sciter.SW_ENABLE_DEBUG,
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("start load html")
	// 加载文件

	if err = w.LoadFile(filepath.Join(path, "demo.html")); err != nil {
		log.Println(err)
	}
	log.Println("load html success")
	// 设置标题
	w.SetTitle("你好，世界")
	// 显示窗口
	w.Show()
	// 运行窗口，进入消息循环
	w.Run()
}
