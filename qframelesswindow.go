package qframelesswindow

import (
	"fmt"
	"runtime"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/svg"
	"github.com/therecipe/qt/widgets"
)

type Edge int

const (
	None        Edge = 0x0
	Left        Edge = 0x1
	Top         Edge = 0x2
	Right       Edge = 0x4
	Bottom      Edge = 0x8
	TopLeft     Edge = 0x10
	TopRight    Edge = 0x20
	BottomLeft  Edge = 0x40
	BottomRight Edge = 0x80
)

type RGB struct {
	R uint16
	G uint16
	B uint16
}

type QToolButtonForNotDarwin struct {
	Widget  *widgets.QWidget
	IconBtn *svg.QSvgWidget
}

type QFramelessWindow struct {
	Window *widgets.QMainWindow
	Widget *widgets.QWidget

	WindowColor *RGB

	borderSize int
	Layout     *widgets.QVBoxLayout

	WindowWidget  *widgets.QFrame
	WindowVLayout *widgets.QVBoxLayout
	shadowMargin  int

	TitleBar          *widgets.QWidget
	TitleBarLayout    *widgets.QHBoxLayout
	TitleLabel        *widgets.QLabel
	TitleBarBtnWidget *widgets.QWidget
	TitleBarBtnLayout *widgets.QHBoxLayout
	TitleColor        *RGB

	// for darwin
	BtnMinimize *widgets.QToolButton
	BtnMaximize *widgets.QToolButton
	BtnRestore  *widgets.QToolButton
	BtnClose    *widgets.QToolButton

	// for windows, linux
	IconMinimize *QToolButtonForNotDarwin
	IconMaximize *QToolButtonForNotDarwin
	IconRestore  *QToolButtonForNotDarwin
	IconClose    *QToolButtonForNotDarwin

	isCursorChanged     bool
	isDragStart         bool
	isLeftButtonPressed bool
	dragPos             *core.QPoint
	pressedEdge         Edge

	Content *widgets.QWidget

	Pos            *core.QPoint
	MousePos       *core.QPoint
	IsMousePressed bool
}

func NewQFramelessWindow() *QFramelessWindow {
	f := &QFramelessWindow{}
	f.Window = widgets.NewQMainWindow(nil, 0)
	f.Widget = widgets.NewQWidget(nil, 0)
	f.SetborderSize(3)
	f.Window.SetCentralWidget(f.Widget)
	f.SetupUI(f.Widget)
	f.SetWindowFlags()
	f.SetWindowShadow(0)
	f.SetAttributes()
	f.SetWindowActions()
	f.SetTitleBarActions()

	return f
}

func (f *QFramelessWindow) SetborderSize(size int) {
	f.borderSize = size
}

func (f *QFramelessWindow) SetWindowShadow(size int) {
	f.shadowMargin = size
	if f.shadowMargin == 0 {
		return
	}
	var alpha int
	if runtime.GOOS == "darwin" {
		alpha = 200
	} else {
		alpha = 100
	}
	shadow := widgets.NewQGraphicsDropShadowEffect(nil)
	shadow.SetBlurRadius((float64)(f.shadowMargin))
	shadow.SetColor(gui.NewQColor3(0, 0, 0, alpha))
	shadow.SetOffset3(-1, 1)
	f.WindowWidget.SetGraphicsEffect(shadow)
	f.Layout.SetContentsMargins(f.shadowMargin, f.shadowMargin, f.shadowMargin, f.shadowMargin)
}

func (f *QFramelessWindow) SetupUI(widget *widgets.QWidget) {
	widget.SetSizePolicy2(widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum, widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum)
	window := f.Window
	window.InstallEventFilter(window)

	f.Layout = widgets.NewQVBoxLayout2(widget)
	f.Layout.SetSpacing(0)

	f.WindowWidget = widgets.NewQFrame(widget, 0)

	f.WindowWidget.SetObjectName("QFramelessWidget")
	f.WindowWidget.SetSizePolicy2(widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum, widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum)
	
	f.Layout.SetContentsMargins(f.shadowMargin, f.shadowMargin, f.shadowMargin, f.shadowMargin)

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
	f.WindowVLayout.SetContentsMargins(0, 0, 0, 0)
	f.WindowVLayout.SetSpacing(0)
	f.WindowWidget.SetLayout(f.WindowVLayout)

	// create titlebar widget
	f.TitleBar = widgets.NewQWidget(f.WindowWidget, 0)
	f.TitleBar.SetObjectName("titleBar")
	f.TitleBar.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Fixed)

	// titleBarLayout is the following structure layout
	// +--+--+--+--+
	// |  |  |  |  |
	// +--+--+--+--+
	f.TitleBarLayout = widgets.NewQHBoxLayout2(f.TitleBar)
	f.TitleBarLayout.SetContentsMargins(0, 0, 0, 0)

	f.TitleLabel = widgets.NewQLabel(nil, 0)
	f.TitleLabel.SetObjectName("TitleLabel")
	f.TitleLabel.SetAlignment(core.Qt__AlignCenter)

	if runtime.GOOS == "darwin" {
		f.SetTitleBarButtonsForDarwin()
	} else {
		f.SetTitleBarButtons()
	}

	// create window content
	f.Content = widgets.NewQWidget(f.WindowWidget, 0)

	// Set widget to layout
	f.WindowVLayout.AddWidget(f.TitleBar, 0, 0)
	f.WindowVLayout.AddWidget(f.Content, 0, 0)

	f.Layout.AddWidget(f.WindowWidget, 0, 0)
}

func (f *QFramelessWindow) SetWidgetColor(red uint16, green uint16, blue uint16, alpha float64) {
	f.WindowColor = &RGB{
		R: red,
		G: green,
		B: blue,
	}
	color := f.WindowColor
	style := fmt.Sprintf("background-color: rgba(%d, %d, %d, %f);", color.R, color.G, color.B, alpha)
	f.Widget.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0); color: rgba(0, 0, 0, 0); }")

	borderSizeString := fmt.Sprintf("%d", f.borderSize*2) + "px"
	f.WindowWidget.SetStyleSheet(fmt.Sprintf(`
	#QFramelessWidget {
		border: 0px solid %s; 
		padding-top: 2px; padding-right: %s; padding-bottom: %s; padding-left: %s; 
		border-radius: %s;
		%s; 
	}`, color.Hex(), borderSizeString, borderSizeString, borderSizeString, borderSizeString, style))
}

func NewQToolButtonForNotDarwin(parent widgets.QWidget_ITF) *QToolButtonForNotDarwin {
	iconSize := 15
	marginTB := iconSize / 6
	marginLR := iconSize / 3

	widget := widgets.NewQWidget(parent, 0)
	widget.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)
	layout := widgets.NewQVBoxLayout2(widget)
	layout.SetContentsMargins(marginLR, marginTB, marginLR, marginTB)
	icon := svg.NewQSvgWidget(nil)
	icon.SetFixedSize2(iconSize, iconSize)

	layout.AddWidget(icon, 0, 0)
	layout.SetAlignment(icon, core.Qt__AlignCenter)

	return &QToolButtonForNotDarwin{
		Widget:  widget,
		IconBtn: icon,
	}
}

func (b *QToolButtonForNotDarwin) SetObjectName(name string) {
	b.IconBtn.SetObjectName(name)
}

func (b *QToolButtonForNotDarwin) Hide() {
	b.Widget.Hide()
}

func (b *QToolButtonForNotDarwin) Show() {
	b.Widget.Show()
}

func (b *QToolButtonForNotDarwin) SetStyle(color *RGB) {
	var backgroundColor string
	if color == nil {
		backgroundColor = "background-color:none;"
	} else {
		backgroundColor = fmt.Sprintf("background-color: rgba(%d, %d, %d, 0.2);", color.R, color.G, color.B)
	}

	b.Widget.SetStyleSheet(fmt.Sprintf(`
	.QWidget { 
		%s;
		border:none;
	}
	`, backgroundColor))
}

func (f *QFramelessWindow) SetTitleBarButtons() {
	iconSize := 15
	f.TitleBarLayout.SetSpacing(1)

	f.IconMinimize = NewQToolButtonForNotDarwin(nil)
	f.IconMinimize.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconMinimize.SetObjectName("IconMinimize")
	f.IconMaximize = NewQToolButtonForNotDarwin(nil)
	f.IconMaximize.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconMaximize.SetObjectName("IconMaximize")
	f.IconRestore = NewQToolButtonForNotDarwin(nil)
	f.IconRestore.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconRestore.SetObjectName("IconRestore")
	f.IconClose = NewQToolButtonForNotDarwin(nil)
	f.IconClose.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconClose.SetObjectName("IconClose")

	f.SetIconsStyle(nil)

	f.IconMinimize.Hide()
	f.IconMaximize.Hide()
	f.IconRestore.Hide()
	f.IconClose.Hide()

	f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignRight)
	f.TitleBarLayout.AddWidget(f.TitleLabel, 0, 0)
	f.TitleBarLayout.AddWidget(f.IconMinimize.Widget, 0, 0)
	f.TitleBarLayout.AddWidget(f.IconMaximize.Widget, 0, 0)
	f.TitleBarLayout.AddWidget(f.IconRestore.Widget, 0, 0)
	f.TitleBarLayout.AddWidget(f.IconClose.Widget, 0, 0)
}

func (f *QFramelessWindow) SetIconsStyle(color *RGB) {
	for _, b := range []*QToolButtonForNotDarwin{
		f.IconMinimize,
		f.IconMaximize,
		f.IconRestore,
		f.IconClose,
	} {
		b.SetStyle(color)
	}
}

func (f *QFramelessWindow) SetTitleBarButtonsForDarwin() {
	btnSizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__ToolButton)
	f.BtnMinimize = widgets.NewQToolButton(f.TitleBar)
	f.BtnMinimize.SetObjectName("BtnMinimize")
	f.BtnMinimize.SetSizePolicy(btnSizePolicy)

	f.BtnMaximize = widgets.NewQToolButton(f.TitleBar)
	f.BtnMaximize.SetObjectName("BtnMaximize")
	f.BtnMaximize.SetSizePolicy(btnSizePolicy)

	f.BtnRestore = widgets.NewQToolButton(f.TitleBar)
	f.BtnRestore.SetObjectName("BtnRestore")
	f.BtnRestore.SetSizePolicy(btnSizePolicy)
	f.BtnRestore.SetVisible(false)

	f.BtnClose = widgets.NewQToolButton(f.TitleBar)
	f.BtnClose.SetObjectName("BtnClose")
	f.BtnClose.SetSizePolicy(btnSizePolicy)

	f.TitleBarLayout.SetSpacing(0)
	f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignLeft)
	f.TitleBarLayout.AddWidget(f.BtnClose, 0, 0)
	f.TitleBarLayout.AddWidget(f.BtnMinimize, 0, 0)
	f.TitleBarLayout.AddWidget(f.BtnMaximize, 0, 0)
	f.TitleBarLayout.AddWidget(f.BtnRestore, 0, 0)
	f.TitleBarLayout.AddWidget(f.TitleLabel, 0, 0)
}

func (f *QFramelessWindow) SetAttributes() {
	f.Window.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.Window.SetAttribute(core.Qt__WA_Hover, true)
	f.Window.SetMouseTracking(true)
}

func (f *QFramelessWindow) SetWindowFlags() {
	f.Window.SetWindowFlag(core.Qt__Window, true)
	f.Window.SetWindowFlag(core.Qt__FramelessWindowHint, true)
	f.Window.SetWindowFlag(core.Qt__NoDropShadowWindowHint, true)
	f.Window.SetWindowFlag(core.Qt__WindowSystemMenuHint, true)
}

func (f *QFramelessWindow) SetTitle(title string) {
	f.TitleLabel.SetText(title)
}

func (f *QFramelessWindow) SetTitleColor(red uint16, green uint16, blue uint16) {
	f.TitleColor = &RGB{
		R: red,
		G: green,
		B: blue,
	}
	f.SetTitleBarColor()
}

func (f *QFramelessWindow) SetTitleBarColor() {
	var color, labelColor *RGB
	window := f.Window
	if window.IsActiveWindow() {
		color = f.TitleColor
	} else {
		color = nil
	}
	labelColor = color
	if labelColor == nil {
		labelColor = &RGB{
			R: 128,
			G: 128,
			B: 128,
		}
	}
	if runtime.GOOS != "darwin" {
		f.TitleLabel.SetStyleSheet(fmt.Sprintf(" *{padding-left: 60px; color: rgb(%d, %d, %d); }", labelColor.R, labelColor.G, labelColor.B))
		f.SetTitleBarColorForNotDarwin(color)
	} else {
		f.TitleLabel.SetStyleSheet(fmt.Sprintf(" *{padding-right: 60px; color: rgb(%d, %d, %d); }", labelColor.R, labelColor.G, labelColor.B))
		f.SetTitleBarColorForDarwin(color)
	}
}

func (f *QFramelessWindow) SetTitleBarColorForNotDarwin(color *RGB) {
	if color == nil {
		color = &RGB{
			R: 128,
			G: 128,
			B: 128,
		}
	}
	SvgMinimize := fmt.Sprintf(`
	<svg style="width:24px;height:24px" viewBox="0 0 24 24">
	<path fill="%s" d="M20,14H4V10H20" />
	</svg>
	`, color.Hex())
	f.IconMinimize.IconBtn.Load2(core.NewQByteArray2(SvgMinimize, len(SvgMinimize)))

	SvgMaximize := fmt.Sprintf(`
	<svg style="width:24px;height:24px" viewBox="0 0 24 24">
	<path fill="%s" d="M4,4H20V20H4V4M6,8V18H18V8H6Z" />
	</svg>
	`, color.Hex())
	f.IconMaximize.IconBtn.Load2(core.NewQByteArray2(SvgMaximize, len(SvgMaximize)))

	SvgRestore := fmt.Sprintf(`
	<svg style="width:24px;height:24px" viewBox="0 0 24 24">
	<path fill="%s" d="M4,8H8V4H20V16H16V20H4V8M16,8V14H18V6H10V8H16M6,12V18H14V12H6Z" />
	</svg>
	`, color.Hex())
	f.IconRestore.IconBtn.Load2(core.NewQByteArray2(SvgRestore, len(SvgRestore)))

	SvgClose := fmt.Sprintf(`
	<svg style="width:24px;height:24px" viewBox="0 0 24 24">
	<path fill="%s" d="M13.46,12L19,17.54V19H17.54L12,13.46L6.46,19H5V17.54L10.54,12L5,6.46V5H6.46L12,10.54L17.54,5H19V6.46L13.46,12Z" />
	</svg>
	`, color.Hex())
	f.IconClose.IconBtn.Load2(core.NewQByteArray2(SvgClose, len(SvgClose)))

	f.IconMinimize.Show()
	f.IconMaximize.Show()
	f.IconRestore.Show()
	f.IconRestore.Widget.SetVisible(false)
	f.IconClose.Show()
}

func (f *QFramelessWindow) SetTitleBarColorForDarwin(color *RGB) {
	var baseStyle, restoreAndMaximizeColor, minimizeColor, closeColor string
	baseStyle = ` #BtnMinimize, #BtnMaximize, #BtnRestore, #BtnClose {
		min-width: 10px;
		min-height: 10px;
		max-width: 10px;
		max-height: 10px;
		border-radius: 6px;
		border-width: 1px;
		border-style: solid;
		margin: 4px;
	}`
	if color != nil {
		restoreAndMaximizeColor = `
			#BtnRestore, #BtnMaximize {
				background-color: rgb(53, 202, 74);
				border-color: rgb(34, 182, 52);
			}
		`
		minimizeColor = `
			#BtnMinimize {
				background-color: rgb(253, 190, 65);
				border-color: rgb(239, 170, 47);
			}
		`
		closeColor = `
			#BtnClose {
				background-color: rgb(252, 98, 93);
				border-color: rgb(239, 75, 71);
			}
		`
	} else {
		restoreAndMaximizeColor = `
			#BtnRestore, #BtnMaximize {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgb(128, 128, 128, 0.2);
			}
		`
		minimizeColor = `
			#BtnMinimize {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgb(128, 128, 128, 0.2);
			}
		`
		closeColor = `
			#BtnClose {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgb(128, 128, 128, 0.2);
			}
		`
	}
	MaximizeColorHover := `
		#BtnMaximize:hover {
			background-color: rgb(53, 202, 74);
			border-color: rgb(34, 182, 52);
			background-image: url(":/icons/MaximizeHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	RestoreColorHover := `
		#BtnRestore:hover {
			background-color: rgb(53, 202, 74);
			border-color: rgb(34, 182, 52);
			background-image: url(":/icons/RestoreHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	minimizeColorHover := `
		#BtnMinimize:hover {
			background-color: rgb(253, 190, 65);
			border-color: rgb(239, 170, 47);
			background-image: url(":/icons/MinimizeHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	closeColorHover := `
		#BtnClose:hover {
			background-color: rgb(252, 98, 93);
			border-color: rgb(239, 75, 71);
			background-image: url(":/icons/CloseHoverDarwin.png");
			background-repeat: no-repeat;
			background-position: center center; 
		}
	`
	f.BtnMinimize.SetStyleSheet(baseStyle + minimizeColor + minimizeColorHover)
	f.BtnMaximize.SetStyleSheet(baseStyle + restoreAndMaximizeColor + MaximizeColorHover)
	f.BtnRestore.SetStyleSheet(baseStyle + restoreAndMaximizeColor + RestoreColorHover)
	f.BtnClose.SetStyleSheet(baseStyle + closeColor + closeColorHover)
}

func (f *QFramelessWindow) SetContent(layout widgets.QLayout_ITF) {
	f.Content.SetLayout(layout)
}

func (f *QFramelessWindow) UpdateWidget() {
	f.Widget.Update()
	f.Window.Update()
}

func (f *QFramelessWindow) SetWindowActions() {
	// Ref: https://stackoverflow.com/questions/5752408/qt-resize-borderless-widget/37507341#37507341
	f.Window.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
		e := gui.NewQMouseEventFromPointer(core.PointerFromQEvent(event))
		switch event.Type() {
		case core.QEvent__ActivationChange:
			f.SetTitleBarColor()

		case core.QEvent__HoverMove:
			f.updateCursorShape(e.GlobalPos())

		case core.QEvent__Leave:
			cursor := gui.NewQCursor()
			cursor.SetShape(core.Qt__ArrowCursor)
			f.Window.SetCursor(cursor)

		case core.QEvent__MouseMove:
			f.mouseMove(e)

		case core.QEvent__MouseButtonPress:
			f.mouseButtonPressed(e)

		case core.QEvent__MouseButtonRelease:
			f.isDragStart = false
			f.isLeftButtonPressed = false
			f.pressedEdge = None

		default:
		}

		return f.Widget.EventFilter(watched, event)
	})
}

func (f *QFramelessWindow) mouseMove(e *gui.QMouseEvent) {
	window := f.Window
	margin := f.shadowMargin

	if f.isLeftButtonPressed {

		if f.pressedEdge != None {

			left := window.FrameGeometry().Left() + margin
			top := window.FrameGeometry().Top() + margin
			right := window.FrameGeometry().Right() - margin
			bottom := window.FrameGeometry().Bottom() - margin

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
			// if newRect.Width() < window.MinimumWidth() {
			// 	left = window.FrameGeometry().X()
			// }
			// if newRect.Height() < window.MinimumHeight() {
			// 	top = window.FrameGeometry().Y()
			// }
			topLeftPoint = core.NewQPoint2(left-margin, top-margin)
			rightBottomPoint = core.NewQPoint2(right+margin, bottom+margin)
			newRect = core.NewQRect2(topLeftPoint, rightBottomPoint)

			window.SetGeometry(newRect)
		}
	}
}

func (f *QFramelessWindow) mouseButtonPressed(e *gui.QMouseEvent) {
	f.pressedEdge = f.calcCursorPos(e.GlobalPos(), f.Window.FrameGeometry())
	if f.pressedEdge != None {
		f.isLeftButtonPressed = true
	}
}

func (f *QFramelessWindow) updateCursorShape(pos *core.QPoint) {
	cursor := gui.NewQCursor()
	if f.Window.IsFullScreen() || f.Window.IsMaximized() {
		if f.isCursorChanged {
			cursor.SetShape(core.Qt__ArrowCursor)
			f.Window.SetCursor(cursor)
		}
	}
	hoverEdge := f.calcCursorPos(pos, f.Window.FrameGeometry())
	f.isCursorChanged = true
	switch hoverEdge {
	case Top, Bottom:
		cursor.SetShape(core.Qt__SizeVerCursor)
		f.Window.SetCursor(cursor)
	case Left, Right:
		cursor.SetShape(core.Qt__SizeHorCursor)
		f.Window.SetCursor(cursor)
	case TopLeft, BottomRight:
		cursor.SetShape(core.Qt__SizeFDiagCursor)
		f.Window.SetCursor(cursor)
	case TopRight, BottomLeft:
		cursor.SetShape(core.Qt__SizeBDiagCursor)
		f.Window.SetCursor(cursor)
	default:
		cursor.SetShape(core.Qt__ArrowCursor)
		f.Window.SetCursor(cursor)
		f.isCursorChanged = false
	}
}

func (f *QFramelessWindow) calcCursorPos(pos *core.QPoint, rect *core.QRect) Edge {
	margin := f.shadowMargin

	doubleBorderSize := f.borderSize * 2
	octupleBorderSize := f.borderSize * 8
	topBorderSize := 2 - 1

	rectX := rect.X() + margin
	rectY := rect.Y() + margin
	rectHeight := rect.Height() - (2 * margin)
	rectWidth := rect.Width() - (2 * margin)

	var onLeft, onRight, onBottom, onTop, onBottomLeft, onBottomRight, onTopRight, onTopLeft bool

	onBottomLeft = (((pos.X() <= (rectX + octupleBorderSize)) && pos.X() >= rectX &&
		(pos.Y() <= (rectY + rectHeight)) && (pos.Y() >= (rectY + rectHeight - doubleBorderSize))) ||
		((pos.X() <= (rectX + doubleBorderSize)) && pos.X() >= rectX &&
			(pos.Y() <= (rectY + rectHeight)) && (pos.Y() >= (rectY + rectHeight - octupleBorderSize))))

	if onBottomLeft {
		return BottomLeft
	}

	onBottomRight = (((pos.X() >= (rectX + rectWidth - octupleBorderSize)) && (pos.X() <= (rectX + rectWidth)) &&
		(pos.Y() >= (rectY + rectHeight - doubleBorderSize)) && (pos.Y() <= (rectY + rectHeight))) ||
		((pos.X() >= (rectX + rectWidth - doubleBorderSize)) && (pos.X() <= (rectX + rectWidth)) &&
			(pos.Y() >= (rectY + rectHeight - octupleBorderSize)) && (pos.Y() <= (rectY + rectHeight))))

	if onBottomRight {
		return BottomRight
	}

	onTopRight = (pos.X() >= (rectX + rectWidth - doubleBorderSize)) && (pos.X() <= (rectX + rectWidth)) &&
		(pos.Y() >= rectY) && (pos.Y() <= (rectY + doubleBorderSize))
	if onTopRight {
		return TopRight
	}

	onTopLeft = pos.X() >= rectX && (pos.X() <= (rectX + doubleBorderSize)) &&
		pos.Y() >= rectY && (pos.Y() <= (rectY + doubleBorderSize))
	if onTopLeft {
		return TopLeft
	}

	onLeft = (pos.X() >= (rectX - doubleBorderSize)) && (pos.X() <= (rectX + doubleBorderSize)) &&
		(pos.Y() <= (rectY + rectHeight - doubleBorderSize)) &&
		(pos.Y() >= rectY+doubleBorderSize)
	if onLeft {
		return Left
	}

	onRight = (pos.X() >= (rectX + rectWidth - doubleBorderSize)) &&
		(pos.X() <= (rectX + rectWidth)) &&
		(pos.Y() >= (rectY + doubleBorderSize)) && (pos.Y() <= (rectY + rectHeight - doubleBorderSize))
	if onRight {
		return Right
	}

	onBottom = (pos.X() >= (rectX + doubleBorderSize)) && (pos.X() <= (rectX + rectWidth - doubleBorderSize)) &&
		(pos.Y() >= (rectY + rectHeight - doubleBorderSize)) && (pos.Y() <= (rectY + rectHeight))
	if onBottom {
		return Bottom
	}

	onTop = (pos.X() >= (rectX + doubleBorderSize)) && (pos.X() <= (rectX + rectWidth - doubleBorderSize)) &&
		(pos.Y() >= rectY) && (pos.Y() <= (rectY + topBorderSize))
	if onTop {
		return Top
	}

	return None
}

func (c *RGB) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", uint8(c.R), uint8(c.G), uint8(c.B))
}
