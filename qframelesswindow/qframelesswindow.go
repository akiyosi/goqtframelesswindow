package qframelesswindow

import (
        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
        "github.com/therecipe/qt/widgets"
)

type QFramelessWindow struct {
	Widget         *widgets.QWidget
	Layout         *widgets.QVBoxLayout

	WindowWidget   *widgets.QWidget
	WindowVLayout  *widgets.QVBoxLayout

	TitleBar       *widgets.QWidget
	TitleBarLayout *widgets.QHBoxLayout
	TitleLabel     *widgets.QLabel
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
	f.setupUI()
	f.setWindowFlags()
	f.setAttribute()

	return f
}

func (f *QFramelessWindow) setupUI() {
	f.Widget = widgets.NewQWidget(nil, 0)
        f.Layout = widgets.NewQVBoxLayout2(f.Widget)
        f.Layout.SetContentsMargins(0, 0, 0, 0)

        f.WindowWidget = widgets.NewQWidget(nil, 0)
        f.WindowWidget.SetObjectName("windowWidget")

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

	// create titlebar widget
	f.TitleBar = widgets.NewQWidget(f.WindowWidget, 0)
        f.TitleBar.SetObjectName("titleBar")
	titlebarSizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__DefaultType)
        f.TitleBar.SetSizePolicy(titlebarSizePolicy)
	f.TitleBar.InstallEventFilter(f.Widget)
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

        f.Layout.AddWidget(f.WindowWidget, 0, 0)
}

func (f *QFramelessWindow) setAttribute() {
	f.Widget.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.Widget.SetAttribute(core.Qt__WA_NoSystemBackground, true)
}

func (f *QFramelessWindow) setWindowFlags() {
	f.Widget.SetWindowFlags(core.Qt__Window | core.Qt__FramelessWindowHint | core.Qt__WindowSystemMenuHint)
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

