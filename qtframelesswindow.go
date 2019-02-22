package qframelesswindow

import (
	"fmt"

        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
        "github.com/therecipe/qt/widgets"
)

type QFramalessWindow {
	window        *widgets.QWindow
	widget        *widgets.QWidget
}

type WindowDragger struct {
	widget          *widgets.QWidget
	pos             *core.QPoint
	mousePos        *core.QPoint
	isMousePressed  bool
}

func NewWindowDragger() *WindowDragger {
	widget := widgets.NewQWidget(nil, 0)
	widget.ConnectMousePressEvent(mousePressEvent)
	widget.ConnectMouseReleaseEvent(mouseReleaseEvent)
	widget.ConnectMoveEvent(moveEvent)
	widget.ConnectMouseDoubleClickEvent(doubleClickEvent)
	widget.ConnectPaintEvent(paintEvent)
}

func (d *WindowDragger) mousePressEvent(event *gui.QMouseEvent) {
	d.isMousePressed = true
	d.mousePos = event.GlobalPos()
}

func (d *WindowDragger) mouseReleaseEvent(event *gui.QMouseEvent) {
	d.isMousePressed = false
}

func (d *WindowDragger) moveEvent(event *gui.QMoveEvent) {
}

func (d *WindowDragger) doubleClickEvent(event *gui.QMouseEvent) {
}

func (d *WindowDragger) paintEvent(event *gui.QPaintEvent) {
}

func NewQFramelessWindow() *FramelessWindow {
	f := &QFramelessWindow{
		window: widgets.NewQMainWindow(nil, 0),
	}
	f.setWindowFlags()
	f.setAttribute()

	f.setupUI()
}

func (f *QFramelessWindow) setupUI() {
	widget := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout(widget)
	layout.SetContentsMargins(0, 0, 0, 0)
}

func (f *QFramelessWindow) setAttribute() {
	f.window.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.window.SetAttribute(core.Qt__WA_NoSystemBackground, true)
}

func (f *QFramelessWindow) setWindowFlags() {
	f.window.SetWindowFlags(core.Qt__FramelessWindowHint, true)
}

