package main

import (
	"fmt"
	
	"github.com/ying32/govcl/vcl"
)

func main() {
	var i = 100_00_0_00
	i++
	fmt.Print(i)
	
	vcl.Application.Initialize()
	mainForm := vcl.Application.CreateForm()
	mainForm.SetCaption("Hello")
	mainForm.EnabledMaximize(false)
	mainForm.ScreenCenter()
	btn := vcl.NewButton(mainForm)
	btn.SetParent(mainForm)
	btn.SetCaption("Hello")
	btn.SetOnClick(func(sender vcl.IObject) {
		vcl.ShowMessage("Hello!")
	})
	vcl.Application.Run()
}