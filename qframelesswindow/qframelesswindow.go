package qframelesswindow

import (
        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
        "github.com/therecipe/qt/widgets"
)

type QFramelessWindow struct {
	Window         *widgets.QMainWindow
	Widget         *widgets.QWidget
	Layout         *widgets.QGridLayout

	WindowWidget   *widgets.QFrame
	WindowVLayout  *widgets.QVBoxLayout

	TitleBar       *widgets.QWidget
	TitleBarLayout *widgets.QHBoxLayout
	TitleLabel     *widgets.QLabel
	TitleBarBtnWidget *widgets.QWidget
	TitleBarBtnLayout *widgets.QHBoxLayout
	BtnMinimize    *widgets.QToolButton
	BtnMaximize    *widgets.QToolButton
	BtnRestore     *widgets.QToolButton
	BtnClose       *widgets.QToolButton

	Content        *widgets.QWidget

	Pos            *core.QPoint
	MousePos       *core.QPoint
	IsMousePressed bool
}

func NewQFramelessWindow() *QFramelessWindow {
	f := &QFramelessWindow{}
	f.Window = widgets.NewQMainWindow(nil, 0)
	f.Widget = widgets.NewQWidget(nil, 0)
	f.Window.SetCentralWidget(f.Widget)
	f.setupUI(f.Widget)
	f.setWindowFlags()
	f.setAttribute()

	return f
}

func (f *QFramelessWindow) setupUI(widget *widgets.QWidget) {
        //f.Layout = widgets.NewQVBoxLayout2(widget)

	widget.SetObjectName("QFramelessWindow")
	widget.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	widget.SetStyleSheet(" .QFramelessWindow { background-color: rgba(0, 0, 0, 0);}")

        f.Layout = widgets.NewQGridLayout(widget)
        f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.Layout.SetSpacing(0)

	// prepare sizegrip
	lefttop := widgets.NewQSizeGrip(widget)
	top := widgets.NewQSizeGrip(widget)
	righttop := widgets.NewQSizeGrip(widget)
	left := widgets.NewQSizeGrip(widget)
	right := widgets.NewQSizeGrip(widget)
	leftbottom := widgets.NewQSizeGrip(widget)
	bottom := widgets.NewQSizeGrip(widget)
	rightbottom := widgets.NewQSizeGrip(widget)

	top.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	left.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
	right.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
	bottom.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)

	lefttop.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); width: 4px; height: 4px;}")
	top.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); height: 4px;}")
	righttop.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); width: 4px; height: 4px;}")
	left.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); width: 4px; }")
	right.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); width: 4px; }")
	leftbottom.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); width: 4px; height: 4px;}")
	bottom.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); height: 4px;}")
	rightbottom.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); width: 4px; height: 4px;}")


        //f.WindowWidget = widgets.NewQWidget(widget, 0)
        f.WindowWidget = widgets.NewQFrame(widget, 0)
        //f.WindowWidget.SetObjectName("QFramelessWidget")
	f.WindowWidget.SetSizePolicy2(widgets.QSizePolicy__Expanding | widgets.QSizePolicy__Maximum , widgets.QSizePolicy__Expanding)
	f.WindowWidget.SetStyleSheet(" .QFrame { border: 2px solid #ccc; border-radius: 11px; background-color: #cccccc; }")

	// windowVLayout is the following structure layout
	// +-----------+
	// |           |
	// +-----------+
	// |           |
	// +-----------+
	// |           |
	// +-----------+
        f.WindowVLayout = widgets.NewQVBoxLayout2(f.WindowWidget)
        f.WindowVLayout.SetContentsMargins(0, 0, 0, 0)
	f.WindowWidget.SetLayout(f.WindowVLayout)

	// create titlebar widget
	f.TitleBar = widgets.NewQWidget(f.WindowWidget, 0)
        f.TitleBar.SetObjectName("titleBar")
	f.TitleBar.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	//f.TitleBar.ConnectEventFilter(f.EventFilter)

	// titleBarLayout is the following structure layout
	// +--+--+--+--+
	// |  |  |  |  |
	// +--+--+--+--+
        f.TitleBarLayout = widgets.NewQHBoxLayout2(f.TitleBar)
        f.TitleBarLayout.SetContentsMargins(0, 0, 0, 0)
        f.TitleBarLayout.SetSpacing(0)

        f.TitleLabel = widgets.NewQLabel(nil, 0)
        f.TitleLabel.SetObjectName("titleLabel")
        f.TitleLabel.SetAlignment(core.Qt__AlignCenter)
        f.TitleBarLayout.AddWidget(f.TitleLabel, 0, 0)

	btnSizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__ToolButton)

	// f.TitleBarBtnWidget = widgets.NewQWidget(f.TitleBar, 0)
	// f.TitleBarBtnWidget.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)
        // f.TitleBarBtnLayout = widgets.NewQHBoxLayout2(f.TitleBarBtnWidget)
        // f.TitleBarBtnLayout.SetContentsMargins(0, 0, 0, 0)
        // f.TitleBarBtnLayout.SetSpacing(0)
        // f.TitleBarLayout.AddWidget(f.TitleBarBtnWidget, 0, 0)

        f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignRight)
        //f.TitleBarBtnLayout.SetAlignment(f.TitleBar, core.Qt__AlignRight)

        f.BtnMinimize = widgets.NewQToolButton(f.TitleBar)
        f.BtnMinimize.SetObjectName("btnMinimize")
        f.BtnMinimize.SetSizePolicy(btnSizePolicy)
        f.TitleBarLayout.AddWidget(f.BtnMinimize, 0, 0)

        f.BtnRestore = widgets.NewQToolButton(f.TitleBar)
        f.BtnRestore.SetObjectName("btnRestore")
        f.BtnRestore.SetSizePolicy(btnSizePolicy)
        f.BtnRestore.SetVisible(false)
        f.TitleBarLayout.AddWidget(f.BtnRestore, 0, 0)

        f.BtnMaximize = widgets.NewQToolButton(f.TitleBar)
        f.BtnMaximize.SetObjectName("btnMaximize")
        f.BtnMaximize.SetSizePolicy(btnSizePolicy)
        f.TitleBarLayout.AddWidget(f.BtnMaximize, 0, 0)

        f.BtnClose = widgets.NewQToolButton(f.TitleBar)
        f.BtnClose.SetObjectName("btnClose")
        f.BtnClose.SetSizePolicy(btnSizePolicy)
        f.TitleBarLayout.AddWidget(f.BtnClose, 0, 0)

	// titleBar connect actions
        f.setTitleBarActions()

	// create window content
        f.Content = widgets.NewQWidget(f.WindowWidget, 0)

	// set widget to layout
        f.WindowVLayout.AddWidget(f.TitleBar, 0, 0)
        f.WindowVLayout.AddWidget(f.Content, 0, 0)


        f.Layout.AddWidget(lefttop, 0, 0, core.Qt__AlignLeft | core.Qt__AlignTop)
        f.Layout.AddWidget(top, 0, 1, core.Qt__AlignTop)
        f.Layout.AddWidget(righttop, 0, 2, core.Qt__AlignRight | core.Qt__AlignTop)
        f.Layout.AddWidget(left, 1, 0, core.Qt__AlignLeft)
        f.Layout.AddWidget(f.WindowWidget, 1, 1, 0)
        f.Layout.AddWidget(right, 1, 2, core.Qt__AlignRight)
        f.Layout.AddWidget(leftbottom, 2, 0, core.Qt__AlignLeft | core.Qt__AlignBottom)
        f.Layout.AddWidget(bottom, 2, 1, core.Qt__AlignBottom)
        f.Layout.AddWidget(rightbottom, 2, 2, core.Qt__AlignRight | core.Qt__AlignBottom)
}

func (f *QFramelessWindow) setAttribute() {
	f.Widget.Window().SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.Widget.Window().SetAttribute(core.Qt__WA_NoSystemBackground, true)
}

func (f *QFramelessWindow) setWindowFlags() {
	//f.Widget.Window().SetWindowFlags(core.Qt__Window | core.Qt__FramelessWindowHint | core.Qt__WindowSystemMenuHint)
	f.Widget.Window().SetWindowFlag(core.Qt__Window, true)
	f.Widget.Window().SetWindowFlag(core.Qt__FramelessWindowHint, true)
	f.Widget.Window().SetWindowFlag(core.Qt__WindowSystemMenuHint, true)
}

func (f *QFramelessWindow) SetTitle(title string) {
	f.TitleLabel.SetText(title)
}

func (f *QFramelessWindow) SetContent(layout widgets.QLayout_ITF) {
	f.Content.SetLayout(layout)
}

func (f *QFramelessWindow) setTitleBarActions() {
	t := f.TitleBar

	// TitleBar Actions
	t.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Raise()
	 	f.IsMousePressed = true
	 	f.MousePos = e.GlobalPos()
		f.Pos = f.Widget.Window().Pos()
	})

	t.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
	 	f.IsMousePressed = false
	})

	t.ConnectMouseMoveEvent(func(e *gui.QMouseEvent) {
		if !f.IsMousePressed {
			return
		}
		x := f.Pos.X() + e.GlobalPos().X() - f.MousePos.X()
		y := f.Pos.Y() + e.GlobalPos().Y() - f.MousePos.Y()
		newPos := core.NewQPoint2(x, y)
		f.Widget.Window().Move(newPos)
	})

	t.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
		if f.BtnMaximize.IsVisible() {
			f.windowMaximize()
		} else {
			f.windowRestore()
		}
	})

	t.ConnectPaintEvent(func(e *gui.QPaintEvent) {
	})

	// Button Actions
	f.BtnMinimize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Window().SetWindowState(core.Qt__WindowMinimized)
	})

	f.BtnMaximize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowMaximize()
	})

	f.BtnRestore.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowRestore()
	})

	f.BtnClose.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
	})
}

func(f *QFramelessWindow) windowMaximize() {
	f.BtnMaximize.SetVisible(false)
	f.BtnRestore.SetVisible(true)
	f.Widget.Window().SetWindowState(core.Qt__WindowMaximized)
}

func(f *QFramelessWindow) windowRestore() {
	f.BtnMaximize.SetVisible(true)
	f.BtnRestore.SetVisible(false)
	f.Widget.Window().SetWindowState(core.Qt__WindowNoState)
}

