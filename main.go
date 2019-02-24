package main

import (
        "github.com/therecipe/qt/widgets"
        "github.com/therecipe/qt/core"
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
	a.win = widgets.NewQMainWindow(nil, 0)
	a.win.SetWindowFlag(core.Qt__Window, true)
	a.win.SetWindowFlag(core.Qt__FramelessWindowHint, true)
	a.win.SetWindowFlag(core.Qt__WindowSystemMenuHint, true)

	a.fw = frameless.NewQFramelessWindow()
	layout := widgets.NewQVBoxLayout()
	a.fw.SetContent(layout)
	a.fw.SetTitle("frameless test")
	a.win.SetCentralWidget(a.fw.Widget)
	// widget := widgets.NewQWidget(nil, 0)
	// layout := widgets.NewQVBoxLayout2(widget)
	// a.win.SetCentralWidget(widget)

	a.win.Show()
	a.fw.Widget.SetFocus2()
	widgets.QApplication_Exec()
}
