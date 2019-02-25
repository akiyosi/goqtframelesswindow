package qframelesswindow

import (
	"fmt"

        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
        "github.com/therecipe/qt/widgets"
)

type QFramelessWindow struct {
	Window         *widgets.QMainWindow
	Widget         *widgets.QWidget
	Layout         *widgets.QGridLayout

	//WindowWidget   *widgets.QFrame
	WindowWidget   *widgets.QWidget
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

	lefttop        *widgets.QSizeGrip
	top            *widgets.QSizeGrip
	righttop       *widgets.QSizeGrip
	left           *widgets.QSizeGrip
	right          *widgets.QSizeGrip
	leftbottom     *widgets.QSizeGrip
	bottom         *widgets.QSizeGrip
	rightbottom    *widgets.QSizeGrip


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
	f.SetStyles("#fff")
	f.setWindowFlags()
	f.setAttribute()

	return f
}

func (f *QFramelessWindow) setupUI(widget *widgets.QWidget) {
        //f.Layout = widgets.NewQVBoxLayout2(widget)

	widget.SetObjectName("QFramelessWindow")
	widget.SetAttribute(core.Qt__WA_TranslucentBackground, true)

        f.Layout = widgets.NewQGridLayout(widget)
        f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.Layout.SetSpacing(0)

        // f.WindowWidget = widgets.NewQFrame(widget, 0)
	// f.WindowWidget.SetStyleSheet(" .QFrame { border: 2px solid #ccc; border-radius: 11px; background-color: #ccc; }")
        f.WindowWidget = widgets.NewQWidget(widget, 0)

        //f.WindowWidget.SetObjectName("QFramelessWidget")
	f.WindowWidget.SetSizePolicy2(widgets.QSizePolicy__Expanding | widgets.QSizePolicy__Maximum , widgets.QSizePolicy__Expanding)

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
        f.TitleLabel.SetObjectName("TitleLabel")
        f.TitleLabel.SetAlignment(core.Qt__AlignCenter)
        f.TitleBarLayout.AddWidget(f.TitleLabel, 0, 0)

	btnSizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__ToolButton)

        f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignRight)

        f.BtnMinimize = widgets.NewQToolButton(f.TitleBar)
        f.BtnMinimize.SetObjectName("BtnMinimize")
        f.BtnMinimize.SetSizePolicy(btnSizePolicy)
        f.TitleBarLayout.AddWidget(f.BtnMinimize, 0, 0)

        f.BtnMaximize = widgets.NewQToolButton(f.TitleBar)
        f.BtnMaximize.SetObjectName("BtnMaximize")
        f.BtnMaximize.SetSizePolicy(btnSizePolicy)
        f.TitleBarLayout.AddWidget(f.BtnMaximize, 0, 0)

        f.BtnRestore = widgets.NewQToolButton(f.TitleBar)
        f.BtnRestore.SetObjectName("BtnRestore")
        f.BtnRestore.SetSizePolicy(btnSizePolicy)
        f.BtnRestore.SetVisible(false)
        f.TitleBarLayout.AddWidget(f.BtnRestore, 0, 0)

        f.BtnClose = widgets.NewQToolButton(f.TitleBar)
        f.BtnClose.SetObjectName("BtnClose")
        f.BtnClose.SetSizePolicy(btnSizePolicy)
        f.TitleBarLayout.AddWidget(f.BtnClose, 0, 0)

	// titleBar connect actions
        f.setTitleBarActions()

	// create window content
        f.Content = widgets.NewQWidget(f.WindowWidget, 0)

	// set widget to layout
        f.WindowVLayout.AddWidget(f.TitleBar, 0, 0)
        f.WindowVLayout.AddWidget(f.Content, 0, 0)

	// prepare sizegrip
	f.lefttop = widgets.NewQSizeGrip(widget)
	f.top = widgets.NewQSizeGrip(widget)
	f.righttop = widgets.NewQSizeGrip(widget)
	f.left = widgets.NewQSizeGrip(widget)
	f.right = widgets.NewQSizeGrip(widget)
	f.leftbottom = widgets.NewQSizeGrip(widget)
	f.bottom = widgets.NewQSizeGrip(widget)
	f.rightbottom = widgets.NewQSizeGrip(widget)

	f.top.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)
	f.left.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
	f.right.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
	f.bottom.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)

        f.Layout.AddWidget(f.lefttop, 0, 0, core.Qt__AlignLeft | core.Qt__AlignTop)
        f.Layout.AddWidget(f.top, 0, 1, core.Qt__AlignTop)
        f.Layout.AddWidget(f.righttop, 0, 2, core.Qt__AlignRight | core.Qt__AlignTop)
        f.Layout.AddWidget(f.left, 1, 0, core.Qt__AlignLeft)
        f.Layout.AddWidget(f.WindowWidget, 1, 1, 0)
        f.Layout.AddWidget(f.right, 1, 2, core.Qt__AlignRight)
        f.Layout.AddWidget(f.leftbottom, 2, 0, core.Qt__AlignLeft | core.Qt__AlignBottom)
        f.Layout.AddWidget(f.bottom, 2, 1, core.Qt__AlignBottom)
        f.Layout.AddWidget(f.rightbottom, 2, 2, core.Qt__AlignRight | core.Qt__AlignBottom)
}

func (f *QFramelessWindow) setAttribute() {
	f.Widget.Window().SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.Widget.Window().SetAttribute(core.Qt__WA_NoSystemBackground, true)
}

func (f *QFramelessWindow) SetStyles(color string) {
	style := fmt.Sprintf("background-color: %s", color)
	f.Widget.SetStyleSheet(fmt.Sprintf(" .QFramelessWindow { border: 2px solid #ccc; border-radius: 11px; %s}", style))
	f.lefttop.SetStyleSheet(fmt.Sprintf("             * { width: 4px; height: 4px; %s}", style))
	f.top.SetStyleSheet(fmt.Sprintf("                 * { height: 4px; %s}", style))
	f.righttop.SetStyleSheet(fmt.Sprintf("            * { width: 4px; height: 4px; %s}", style))
	f.left.SetStyleSheet(fmt.Sprintf("                * { width: 4px;  %s}", style))
	f.right.SetStyleSheet(fmt.Sprintf("               * { width: 4px;  %s}", style))
	f.leftbottom.SetStyleSheet(fmt.Sprintf("          * { width: 4px; height: 4px; %s}", style))
	f.bottom.SetStyleSheet(fmt.Sprintf("              * { height: 4px; %s}", style))
	f.rightbottom.SetStyleSheet(fmt.Sprintf("         * { width: 4px; height: 4px; %s}", style))
	f.WindowWidget.SetStyleSheet(fmt.Sprintf(" .QWidget { %s}", style))

	f.BtnMinimize.SetStyleSheet(`
	#BtnMinimize { 
		background-color:none;
		border:none;
		background-image: url(":/icons/Minimize.png"); 
		background-repeat: no-repeat; 
		background-position: center center; 
	}
	#BtnMinimize:hover { 
		background-color:none;
		border:none;
		background-image: url(":/icons/MinimizeHover.png"); 
		background-repeat: no-repeat; 
		background-position: center center;
	}
	`)
	f.BtnMaximize.SetStyleSheet(`
	#BtnMaximize { 
		background-color:none;
		border:none;
		background-image: url(":/icons/Maximize.png");
		background-repeat: no-repeat; 
		background-position: center center; 
	}
	#BtnMaximize:hover { 
		background-color:none;
		border:none;
		background-image: url(":/icons/MaximizeHover.png"); 
		background-repeat: no-repeat; 
		background-position: center center; 
	}
	`)

	f.BtnRestore.SetStyleSheet(`
	#BtnRestore { 
		background-color:none;
		border:none;
		background-image: url(":/icons/Restore.png");
		background-repeat: no-repeat; 
		background-position: center center; 
	}
	#BtnRestore:hover { 
		background-color:none;
		border:none;
		background-image: url(":/icons/RestoreHover.png"); 
		background-repeat: no-repeat; 
		background-position: center center; 
	}
	`)

	f.BtnClose.SetStyleSheet(`
	#BtnClose { 
		background-color:none;
		border:none;
		background-image: url(":/icons/Close.png");
		background-repeat: no-repeat; 
		background-position: center center; 
	}
	#BtnClose:hover { 
		background-color:none;
		border:none;
		background-image: url(":/icons/CloseHover.png");
		background-repeat: no-repeat;
		background-position: center center; 
	}
	`)
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

