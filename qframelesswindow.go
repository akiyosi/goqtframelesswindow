package qframelesswindow

import (
	"fmt"

        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
        "github.com/therecipe/qt/widgets"
)

type QFramalessWindow struct {
	widget         *widgets.QWidget
	layout         *widgets.QVBoxLayout

	windowWidget   *widgets.QWidget
	windowVLayout  *widgets.QVBoxLayout

	titleBar       *widgets.QWidget
	titleBarLayout *widgets.QHBoxLayout
	titleLabel     *widgets.QLabel
	btnMinimize    *widgets.QToolButton
	btnMaximize    *widgets.QToolButton
	btnRestore     *widgets.QToolButton
	btnClose       *widgets.QToolButton

	content        *widgets.QWidget

	pos            *core.QPoint
	mousePos       *core.QPoint
	isMousePressed bool
}


// func (d *WindowDragger) mousePressEvent(event *gui.QMouseEvent) {
// 	d.isMousePressed = true
// 	d.mousePos = event.GlobalPos()
// }
// 
// func (d *WindowDragger) mouseReleaseEvent(event *gui.QMouseEvent) {
// 	d.isMousePressed = false
// }
// 
// func (d *WindowDragger) moveEvent(event *gui.QMoveEvent) {
// }
// 
// func (d *WindowDragger) doubleClickEvent(event *gui.QMouseEvent) {
// }
// 
// func (d *WindowDragger) paintEvent(event *gui.QPaintEvent) {
// }

func NewQFramelessWindow() *FramelessWindow {
	f := &QFramelessWindow{}
	f.setupUI()
	f.setWindowFlags()
	f.setAttribute()
}

func (f *QFramelessWindow) setupUI() {
	f.widget = widgets.NewQWidget(nil, 0)
        f.layout = widgets.NewQVBoxLayout2(f.widget)
        f.layout.setContentsMargins(0, 0, 0, 0)


        f.windowWidget = widgets.NewQWidget(nil, 0)
        f.windowWidget.setObjectName("windowWidget")

	// windowVLayout is the following structure layout
	// +-----------+
	// |           |
	// +-----------+
	// |           |
	// +-----------+
	// |           |
	// +-----------+
        f.windowVLayout = widgets.NewQVBoxLayout2(f.windowWidget)
        f.windowVLayout.setContentsMargins(0, 0, 0, 0)

	// create titlebar widget
	f.titleBar = widgets.NewQWidget(f.windowWidget, 0)
        f.titleBar.setObjectName("titleBar")
        f.titleBar.setSizePolicy(widgets.QSizePolicy(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Fixed))

	// titleBarLayout is the following structure layout
	// +--+--+--+--+
	// |  |  |  |  |
	// +--+--+--+--+
        f.titleBarLayout = widgets.NewQHBoxLayout2(f.titleBar.widget)
        f.titleBarLayout.setContentsMargins(0, 0, 0, 0)
        f.titleBarLayout.setSpacing(0)

        f.titleLabel = QLabel("Title")
        f.titleLabel.setObjectName("titleLabel")
        f.titleLabel.setAlignment(Qt.AlignCenter)
        f.titleBarLayout.addWidget(f.titleLabel)

        btnSizePolicy = QSizePolicy(QSizePolicy.Fixed, QSizePolicy.Fixed)

        f.btnMinimize = widgets.NewQToolButton(f.titleBar.widget)
        f.btnMinimize.setObjectName("btnMinimize")
        f.btnMinimize.setSizePolicy(btnSizePolicy)
        f.titleBarLayout.addWidget(f.btnMinimize)

        f.btnRestore = widgets.NewQToolButton(f.titleBar.widget)
        f.btnRestore.setObjectName("btnRestore")
        f.btnRestore.setSizePolicy(btnSizePolicy)
        f.btnRestore.setVisible(False)
        f.titleBarLayout.addWidget(f.btnRestore)

        f.btnMaximize = widgets.NewQToolButton(f.titleBar.widget)
        f.btnMaximize.setObjectName("btnMaximize")
        f.btnMaximize.setSizePolicy(btnSizePolicy)
        f.titleBarLayout.addWidget(f.btnMaximize)

        f.btnClose = widgets.NewQToolButton(f.titleBar.widget)
        f.btnClose.setObjectName("btnClose")
        f.btnClose.setSizePolicy(btnSizePolicy)
        f.titleBarLayout.addWidget(f.btnClose)

	// titleBar connect actions
        f.connectTitleBarActions(titleBar, f.windowWidget)

	// create window content
        f.content = widgets.NewQWidget(f.windowWidget, 0)

	// set widget to layout
        f.windowVLayout.addWidget(f.titleBar.widget)
        f.windowVLayout.addWidget(f.content)

        f.layout.addWidget(f.windowWidget)
}

func (f *QFramelessWindow) setAttribute() {
	f.widget.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.widget.SetAttribute(core.Qt__WA_NoSystemBackground, true)
}

func (f *QFramelessWindow) setWindowFlags() {
	f.widget.SetWindowFlags(core.Qt__Window, true)
	f.widget.SetWindowFlags(core.Qt__FramelessWindowHint, true)
	f.widget.SetWindowFlags(core.Qt__WindowSystemMenuHint, true)
}

func (f *QFramelessWindow) setTitleTitle(title string) {
	f.titleLabel.SetText(title)
}


func (f *QFramelessWindow) connectTitleBarActions(w widgets.QWidget, parent widgets.Qwidget) {
	t := f.titleBar

	t.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
	 	f.isMousePressed = true
	 	f.mousePos = event.GlobalPos()
	})

	t.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
	 	f.isMousePressed = false
	})

	t.ConnectMoveEvent(func(e *gui.QMoveEvent) {
	})

	t.ConnectMouseDoubleClickEvent(func(e *gui.QmouseEvent) {
	})

	t.ConnectPaintEvent(func(e *gui.QPaintEvent) {
	})

}
