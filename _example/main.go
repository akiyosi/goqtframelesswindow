// package main
//
// import (
// 	"time"
//
// 	frameless "github.com/akiyosi/goqtframelesswindow"
// 	"github.com/therecipe/qt/widgets"
// )
//
// type framelessTest struct {
// 	fw *frameless.QFramelessWindow
//
// 	app *widgets.QApplication
// 	win *widgets.QMainWindow
// }
//
// func main() {
// 	a := &framelessTest{}
// 	a.app = widgets.NewQApplication(0, nil)
//
// 	// a.fw = frameless.CreateQFramelessWindow(1.0, false)
// 	a.fw = frameless.CreateQFramelessWindow(1.0)
// 	layout := widgets.NewQVBoxLayout()
// 	a.fw.SetupContent(layout)
// 	a.fw.SetupWidgetColor(30, 30, 30)
// 	a.fw.SetupTitle("frameless test")
// 	a.fw.SetupTitleColor(220, 220, 220)
//
// 	label := widgets.NewQLabel(nil, 0)
// 	label.SetStyleSheet(" * { color: #eee; }")
// 	label.SetText(`	General relativity (GR, also known as the general theory of relativity or GTR) is
// 	the geometric theory of gravitation published by Albert Einstein in 1915 and
// 	the current description of gravitation in modern physics. `)
//
// 	go func() {
// 		time.Sleep(13000 * time.Millisecond)
// 		label.SetText(`	Update! `)
// 		time.Sleep(3000 * time.Millisecond)
// 		label.SetText(`	Update 2! `)
// 	}()
//
// 	layout.AddWidget(label, 0, 0)
//
// 	a.fw.Show()
// 	a.fw.Widget.SetFocus2()
// 	// a.fw.WindowMaximize()
// 	a.app.Exec()
//
// }
