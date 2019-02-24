package main

import (
	"fmt"

	. "github.com/akiyosi/qframelesswindow"
        "github.com/therecipe/qt/widgets"
)


func main() {
	app = widgets.NewQApplication(0, nil)
	window = widgets.NewQMainWindow(nil, 0)

	fw := NewQFramelessWindow()
        layout = widgets.NewQVBoxLayout()
	fw.SetContent(layout)

	window.SetCentralWidget(fw.widget)
	window.Show()
	fw.widget.SetFocus2()
	widgets.QApplication_Exec()
}
