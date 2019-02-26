package qframelesswindow

import (
	"fmt"

        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
        "github.com/therecipe/qt/widgets"
)

type Edge int

const (
        None Edge = 0x0
        Left Edge = 0x1
        Top Edge = 0x2
        Right Edge = 0x4
        Bottom Edge = 0x8
        TopLeft Edge = 0x10
        TopRight Edge = 0x20
        BottomLeft Edge = 0x40
        BottomRight Edge = 0x80
)

type QFramelessWindow struct {
	Window         *widgets.QMainWindow
	Widget         *widgets.QWidget

	borderSize     int
	Layout         *widgets.QVBoxLayout
	// Layout         *widgets.QGridLayout

	WindowWidget   *widgets.QFrame
	// WindowWidget   *widgets.QWidget
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

	isCursorChanged bool
	isDragStart     bool
	dragPos         *core.QPoint
	pressedEdge     Edge

	Content        *widgets.QWidget

	Pos            *core.QPoint
	MousePos       *core.QPoint
	IsMousePressed bool
}

func NewQFramelessWindow() *QFramelessWindow {
	f := &QFramelessWindow{}
	f.Window = widgets.NewQMainWindow(nil, 0)
	f.Widget = widgets.NewQWidget(nil, 0)
	f.setborderSize(4)
	f.Window.SetCentralWidget(f.Widget)
	f.setupUI(f.Widget)
	f.setWindowFlags()
	f.setAttribute()
        f.setTitleBarActions()
	f.setWindowActions()

	return f
}

func (f *QFramelessWindow) setborderSize(size int) {
	f.borderSize = size
}

func (f *QFramelessWindow) setupUI(widget *widgets.QWidget) {
        //f.Layout = widgets.NewQVBoxLayout2(widget)

	widget.SetObjectName("QFramelessWindow")
	window := widget.Window()
	window.InstallEventFilter(window)

        // f.Layout = widgets.NewQGridLayout(widget)
        f.Layout = widgets.NewQVBoxLayout2(widget)
        f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.Layout.SetSpacing(0)

        f.WindowWidget = widgets.NewQFrame(widget, 0)
        // f.WindowWidget = widgets.NewQWidget(widget, 0)

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
        f.WindowVLayout.SetContentsMargins(f.borderSize, f.borderSize, f.borderSize, 0)
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


	// create window content
        f.Content = widgets.NewQWidget(f.WindowWidget, 0)

	// set widget to layout
        f.WindowVLayout.AddWidget(f.TitleBar, 0, 0)
        f.WindowVLayout.AddWidget(f.Content, 0, 0)

        f.Layout.AddWidget(f.WindowWidget, 0, 0)
}

func (f *QFramelessWindow) setAttribute() {
	f.Widget.Window().SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.Widget.Window().SetAttribute(core.Qt__WA_NoSystemBackground, true)
	f.Widget.Window().SetAttribute(core.Qt__WA_Hover, true)
	f.Widget.Window().SetMouseTracking(true)
	f.Widget.SetAttribute(core.Qt__WA_TranslucentBackground, true)
}

func (f *QFramelessWindow) SetStyles(color string) {
	style := fmt.Sprintf("background-color: %s", color)
	// f.Widget.SetStyleSheet(fmt.Sprintf(" .QFramelessWindow { border: 2px solid #ccc; border-radius: 11px; %s}", style))
	f.Widget.SetStyleSheet("* { background-color: rgba(0, 0, 0, 0); }")

	// f.WindowWidget.SetStyleSheet(fmt.Sprintf(" .QWidget { %s}", style))
	f.WindowWidget.SetStyleSheet(fmt.Sprintf(" .QFrame { border: 2px solid %s; border-radius: 6px; %s }", color, style))

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

func (f *QFramelessWindow) setWindowActions() {
	// Ref: https://stackoverflow.com/questions/5752408/qt-resize-borderless-widget/37507341#37507341
	f.Widget.Window().ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
		e := gui.NewQMouseEventFromPointer(core.PointerFromQEvent(event))
		switch event.Type() {
		case core.QEvent__HoverMove :
	 		f.updateCursorShape(e.GlobalPos())

		case core.QEvent__Leave :
			f.Widget.Window().UnsetCursor()

		case core.QEvent__MouseMove :
			if f.isDragStart {
				startPos := f.Widget.Window().FrameGeometry().TopLeft()
				newX :=startPos.X() + e.Pos().X() - f.dragPos.X()
				newY :=startPos.Y() + e.Pos().Y() - f.dragPos.Y()
				newPoint := core.NewQPoint2(newX, newY)
				f.Widget.Window().Move(newPoint)
			}
			if f.pressedEdge != None {

				left := f.Widget.Window().FrameGeometry().Left()
				top := f.Widget.Window().FrameGeometry().Top()
				right := f.Widget.Window().FrameGeometry().Right()
				bottom := f.Widget.Window().FrameGeometry().Bottom()

				switch f.pressedEdge {
				case Top:
					top = e.GlobalPos().Y()
				case Bottom:
					bottom = e.GlobalPos().Y()
				case Left:
					left = e.GlobalPos().X()
				case Right:
					right = e.GlobalPos().X()
				case TopLeft:
					top = e.GlobalPos().Y()
					left = e.GlobalPos().X()
				case TopRight:
					top = e.GlobalPos().Y()
					right = e.GlobalPos().X()
				case BottomLeft:
					bottom = e.GlobalPos().Y()
					left = e.GlobalPos().X()
				case BottomRight:
					bottom = e.GlobalPos().Y()
					right = e.GlobalPos().X()
				default:
				}

				topLeftPoint := core.NewQPoint2(left, top)
				rightBottomPoint := core.NewQPoint2(right, bottom)
				newRect := core.NewQRect2(topLeftPoint, rightBottomPoint)
				if newRect.Width() < f.Widget.Window().MinimumWidth() {
					left = f.Widget.Window().FrameGeometry().X()
				}
				if newRect.Height() < f.Widget.Window().MinimumHeight() {
					top = f.Widget.Window().FrameGeometry().Y()
				}
				topLeftPoint = core.NewQPoint2(left, top)
				rightBottomPoint = core.NewQPoint2(right, bottom)
				newRect = core.NewQRect2(topLeftPoint, rightBottomPoint)

				f.Widget.Window().SetGeometry(newRect)
			}
		case core.QEvent__MouseButtonPress :
			f.pressedEdge = f.calcCursorPos(e.GlobalPos(), f.Widget.Window().FrameGeometry())
			if f.pressedEdge != None {
				margins := core.NewQMargins2(f.borderSize, f.borderSize, f.borderSize, f.borderSize)
				if f.Widget.Window().Rect().MarginsRemoved(margins).Contains3(e.Pos().X(), e.Pos().Y()) {
					f.isDragStart = true
					f.dragPos = e.Pos()
				}
			}
		case core.QEvent__MouseButtonRelease :
			f.isDragStart = false
			f.pressedEdge  = None

		default:
		}

		return f.Widget.EventFilter(watched, event)
	})
}

func (f *QFramelessWindow) updateCursorShape(pos *core.QPoint) {
	if f.Widget.Window().IsFullScreen() || f.Widget.Window().IsMaximized() {
		if f.isCursorChanged {
			f.Widget.Window().UnsetCursor()
		}
	}
	hoverEdge := f.calcCursorPos(pos, f.Widget.Window().FrameGeometry())
	f.isCursorChanged = true
	cursor := gui.NewQCursor()
	switch hoverEdge {
	case Top, Bottom:
		cursor.SetShape(core.Qt__SizeVerCursor)
		f.Widget.Window().SetCursor(cursor)
	case Left, Right:
		cursor.SetShape(core.Qt__SizeHorCursor)
		f.Widget.Window().SetCursor(cursor)
	case TopLeft, BottomRight:
		cursor.SetShape(core.Qt__SizeFDiagCursor)
		f.Widget.Window().SetCursor(cursor)
	case TopRight, BottomLeft:
		cursor.SetShape(core.Qt__SizeBDiagCursor)
		f.Widget.Window().SetCursor(cursor)
	default:
		f.Widget.Window().UnsetCursor()
		f.isCursorChanged = false
	}
}

func (f *QFramelessWindow) calcCursorPos(pos *core.QPoint, rect *core.QRect) Edge {
	var onLeft, onRight, onBottom, onTop, onBottomLeft, onBottomRight, onTopRight, onTopLeft bool
	onLeft = (pos.X() >= (rect.X() - f.borderSize)) && (pos.X() <= (rect.X() + f.borderSize)) &&
	         (pos.Y() <= (rect.Y() + rect.Height() - f.borderSize)) &&
		 (pos.Y() >= rect.Y() + f.borderSize)
	if onLeft {
		return Left
	}

	onRight = (pos.X() >= (rect.X() + rect.Width() - f.borderSize)) &&
	          (pos.X() <= (rect.X() + rect.Width())) &&
		  (pos.Y() >= (rect.Y() + f.borderSize)) && (pos.Y() <= (rect.Y() + rect.Height() - f.borderSize))
	if onRight {
		return Right
	}
	
	onBottom = (pos.X() >= (rect.X() + f.borderSize)) && (pos.X() <= (rect.X() + rect.Width() - f.borderSize)) &&
	           (pos.Y() >= (rect.Y() + rect.Height() - f.borderSize)) && (pos.Y() <= (rect.Y() + rect.Height()))
	if onBottom {
		return Bottom
	}

	onTop = (pos.X() >= (rect.X() + f.borderSize)) && (pos.X() <= (rect.X() + rect.Width() - f.borderSize)) &&
		(pos.Y() >= rect.Y()) && (pos.Y() <= (rect.Y() + f.borderSize))
	if onTop {
		return Top
	}
	
	onBottomLeft = (pos.X() <= (rect.X() + f.borderSize)) && pos.X() >= rect.X() &&
		       (pos.Y() <= (rect.Y() + rect.Height())) && (pos.Y() >= (rect.Y() + rect.Height() - f.borderSize))
	if onBottomLeft {
		return BottomLeft
	}
	
	onBottomRight = (pos.X() >= (rect.X() + rect.Width() - f.borderSize)) && (pos.X() <= (rect.X() + rect.Width())) &&
	                (pos.Y() >= (rect.Y() + rect.Height() - f.borderSize)) && (pos.Y() <= (rect.Y() + rect.Height()))
	if onBottomRight {
		return BottomRight
	}
	
	onTopRight = (pos.X() >= (rect.X() + rect.Width() - f.borderSize)) && (pos.X() <= (rect.X() + rect.Width())) &&
		     (pos.Y() >= rect.Y()) && (pos.Y() <= (rect.Y() + f.borderSize))
	if onTopRight {
		return TopRight
	}
	
	onTopLeft = pos.X() >= rect.X() && (pos.X() <= (rect.X() + f.borderSize)) &&
		    pos.Y() >= rect.Y() && (pos.Y() <= (rect.Y() + f.borderSize))
	if onTopLeft {
		return TopLeft
	}
	
	return None
}
