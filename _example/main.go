package main

import (
	"time"

	frameless "github.com/akiyosi/goqtframelesswindow"
	"github.com/therecipe/qt/widgets"
)

type framelessTest struct {
	fw *frameless.QFramelessWindow

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
	a.fw.SetWidgetColor(30, 30, 30, 0.7)
	a.fw.SetTitle("frameless test")
	a.fw.SetTitleColor(200, 200, 200)

	label := widgets.NewQLabel(nil, 0)
	label.SetStyleSheet(" * { color: #eee; }")
	label.SetText(`	General relativity (GR, also known as the general theory of relativity or GTR) is
	the geometric theory of gravitation published by Albert Einstein in 1915 and 
	the current description of gravitation in modern physics. `)

	go func() {
		time.Sleep(13000 * time.Millisecond)
		label.SetText(`	Update! `)
		time.Sleep(3000 * time.Millisecond)
		label.SetText(`	Update 2! `)
	}()

	layout.AddWidget(label, 0, 0)

	// In Windows, signal arrived during external code execution
	// In MacOS, bad access
	a.fw.SetNativeEvent(a.app)

	a.win.Show()
	a.fw.Widget.SetFocus2()
	a.app.Exec()

}
