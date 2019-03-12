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
	a.fw.SetWidgetColor(30, 30, 30, 0.5)
	a.fw.SetTitle("frameless test")
	a.fw.SetTitleColor(200, 200, 200)

	label := widgets.NewQLabel(nil, 0)
	layout.AddWidget(label, 0, 0)
	label.SetStyleSheet(" * { color: #eee; }")
	label.SetText(`
	General relativity (GR, also known as the general theory of relativity or GTR) is
	the geometric theory of gravitation published by Albert Einstein in 1915 and 
	the current description of gravitation in modern physics. `)

	go func() {
		time.Sleep(16 * 1000 * time.Millisecond)
		label.SetText("update text 1")

		time.Sleep(5 * 1000 * time.Millisecond)
		label.SetText(`
		Loop quantum gravity (LQG) is a theory of quantum gravity, 
		merging quantum mechanics and general relativity, making it
		a possible candidate for a theory of everything.
		Its goal is to unify gravity in a common theoretical framework 
		with the other three fundamental forces of nature, beginning with
		relativity and adding quantum features. It competes with
		string theory that begins with quantum field theory and adds gravity.
		`)
	}()

	//a.fw.SetNativeEvent(a.app)

	a.win.Show()
	a.fw.Widget.SetFocus2()
	widgets.QApplication_Exec()
}
