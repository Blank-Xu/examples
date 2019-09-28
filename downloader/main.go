package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/andlabs/ui"
)

func main() {
	ui1 := ui.NewWindow("")
	
	app := app.New()
	
	w := app.NewWindow("hello")
	w.SetContent(
		widget.NewVBox(
			widget.NewLabel("hello fyne"),
			widget.NewButton("quit", func() {
				app.Quit()
			}),
			))
	w.ShowAndRun()
}