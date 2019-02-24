package main

import (
        "github.com/therecipe/qt/widgets"
	frameless "github.com/akiyosi/goqtframelesswindow/qframelesswindow"
)

type framelessTest struct {
	fw  *frameless.QFramelessWindow

	app *widgets.QApplication
	win *widgets.QMainWindow
}

func main() {
	a := &framelessTest{}
	a.app = widgets.NewQApplication(0, nil)

	a.fw = frameless.NewQFramelessWindow()
	a.win = a.fw.Window
	layout := widgets.NewQVBoxLayout()
	a.fw.SetContent(layout)
	a.fw.SetTitle("frameless test")

	a.win.Show()
	a.fw.Widget.SetFocus2()
	widgets.QApplication_Exec()
}
