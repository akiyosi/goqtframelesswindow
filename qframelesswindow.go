package qframelesswindow

import (
	"fmt"
	"runtime"
	"time"

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
	f       *QFramelessWindow
	Widget  *widgets.QWidget
	IconBtn *svg.QSvgWidget
	isHover bool
}

type QFramelessWindow struct {
	widgets.QMainWindow

	IsBorderless bool

	WindowColor      *RGB
	WindowColorAlpha float64

	Widget  *widgets.QWidget
	Layout  *widgets.QVBoxLayout
	Content *widgets.QWidget

	WindowWidget  *widgets.QFrame
	WindowVLayout *widgets.QVBoxLayout
	borderSize    int
	windowgap     int
	minimumWidth  int
	minimumHeight int

	TitleBar       *widgets.QWidget
	TitleBarLayout *widgets.QHBoxLayout
	// TitleLabel        *widgets.QLabel
	TitleIconLabel    *widgets.QLabel
	TitleStringLabel  *widgets.QLabel
	TitleLabel        *widgets.QWidget
	TitleBarBtnWidget *widgets.QWidget
	TitleBarBtnLayout *widgets.QHBoxLayout
	TitleColor        *RGB
	TitleBarMousePos  *core.QPoint
	IsTitlebarHidden  bool
	IsTitleBarPressed bool
	IsTitleIconShown  bool

	// for windows, linux
	IconMinimize *QToolButtonForNotDarwin
	IconMaximize *QToolButtonForNotDarwin
	IconRestore  *QToolButtonForNotDarwin
	IconClose    *QToolButtonForNotDarwin

	isCursorChanged     bool
	isDragStart         bool
	isLeftButtonPressed bool
	dragPos             *core.QPoint
	hoverEdge           Edge
	Position            *core.QPoint
	MousePos            [2]int

	borderless bool
}

func CreateQFramelessWindow(a ...interface{}) *QFramelessWindow {
	alpha := 1.0
	isBorderless := true
	for _, vITF := range a {
		switch vITF.(type) {
		case float64:
			v := vITF.(float64)
			if v >= 0.0 && v <= 1.0 {
				alpha = v
			}
		case bool:
			isBorderless = vITF.(bool)
		default:
		}
	}

	f := NewQFramelessWindow(nil, 0)
	// f.SetupBorderSize(2)
	f.WindowColor = &RGB{255, 255, 255}
	f.WindowColorAlpha = alpha
	f.IsBorderless = isBorderless

	// for windows
	// if f.WindowColorAlpha == 1.0 && f.IsBorderless {
	// 	f.SetupNativeEvent()
	// } else if f.WindowColorAlpha < 1.0 && f.IsBorderless {
	// 	f.SetupNativeEvent2()
	// }
	if f.IsBorderless {
		f.SetupNativeEvent()
	}

	f.Widget = widgets.NewQWidget(nil, 0)
	f.Widget.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0.0); color: rgba(0, 0, 0, 0); }")
	f.SetCentralWidget(f.Widget)
	f.SetupUI(f.Widget)
	f.SetupMinimumSize(400, 300)

	// if isBorderless is false, then we return normal qmainwindow
	if !f.IsBorderless {
		f.HideTitlebar()
		return f
	}
	if f.IsTitlebarHidden {
		f.TitleBar.Hide()
	}

	if f.WindowColorAlpha < 1.0 {
		f.SetupWindowFlags()
		f.SetupAttributes()
	}

	f.SetupWindowActions()
	f.SetupTitleBarActions()
	f.ShowButtons()

	return f
}

func (f *QFramelessWindow) SetupMinimumSize(w int, h int) {
	f.SetMinimumSize2(w, h)
	f.Widget.SetMinimumSize2(w, h)
	f.WindowWidget.SetMinimumSize2(w, h)
	f.minimumWidth = w
	f.minimumHeight = h
}

func (f *QFramelessWindow) BorderSize() int {
	return f.borderSize
}

func (f *QFramelessWindow) SetupBorderSize(size int) {
	f.borderSize = size
}

func (f *QFramelessWindow) WindowGap() int {
	return f.windowgap
}

func (f *QFramelessWindow) SetupWindowGap(size int) {
	f.windowgap = size
}

func (f *QFramelessWindow) SetupUI(widget *widgets.QWidget) {
	if f.IsBorderless {
		f.InstallEventFilter(f)
	}

	widget.SetSizePolicy2(widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum, widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum)
	f.Layout = widgets.NewQVBoxLayout2(widget)
	f.Layout.SetSpacing(0)
	f.Layout.SetContentsMargins(0, 0, 0, 0)

	f.WindowWidget = widgets.NewQFrame(widget, 0)
	f.WindowWidget.SetObjectName("QFramelessWidget")
	f.WindowWidget.SetSizePolicy2(widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum, widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum)

	f.Layout.AddWidget(f.WindowWidget, 0, 0)

	if !f.IsBorderless {
		return
	}
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
	f.WindowVLayout.SetSpacing(0)

	f.newTitlebar()

	// create window content
	f.Content = widgets.NewQWidget(f.WindowWidget, 0)

	// Set widget to layout
	f.WindowVLayout.AddWidget(f.TitleBar, 0, 0)
	f.WindowVLayout.AddWidget(f.Content, 0, 0)
	f.WindowWidget.SetLayout(f.WindowVLayout)
}

func (f *QFramelessWindow) newTitlebar() {
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

	f.TitleBarBtnWidget = widgets.NewQWidget(nil, 0)
	f.TitleBarBtnLayout = widgets.NewQHBoxLayout2(nil)
	f.TitleBarBtnLayout.SetContentsMargins(0, 0, 0, 0)
	f.TitleBarBtnLayout.SetSpacing(0)
	f.TitleBarBtnWidget.SetLayout(f.TitleBarBtnLayout)

	// f.TitleLabel = widgets.NewQLabel(nil, 0)
	// f.TitleLabel.SetObjectName("TitleLabel")
	// f.TitleLabel.SetAlignment(core.Qt__AlignCenter)

	f.TitleLabel = widgets.NewQWidget(nil, 0)
	titleLabelLayout := widgets.NewQHBoxLayout2(nil)
	titleLabelLayout.SetContentsMargins(0, 0, 0, 0)
	titleLabelLayout.SetSpacing(int(float64(f.borderSize) * 2.5))
	f.TitleLabel.SetLayout(titleLabelLayout)

	f.TitleIconLabel = widgets.NewQLabel(nil, 0)
	f.TitleIconLabel.SetContentsMargins(0, 2, 0, 0)
	f.TitleStringLabel = widgets.NewQLabel(nil, 0)
	f.TitleStringLabel.SetContentsMargins(0, 0, 0, 0)

	titleLabelLayout.AddWidget(f.TitleIconLabel, 0, 0)
	titleLabelLayout.AddWidget(f.TitleStringLabel, 0, 0)

	titleLabelLayout.SetAlignment(f.TitleIconLabel, core.Qt__AlignCenter)
	titleLabelLayout.SetAlignment(f.TitleStringLabel, core.Qt__AlignCenter)

	f.TitleIconLabel.Hide()

	if runtime.GOOS == "darwin" {
		f.SetTitleBarButtonsForDarwin()
	} else {
		f.SetTitleBarButtons()
	}
}

func (f *QFramelessWindow) SetupWidgetColor(red uint16, green uint16, blue uint16) {
	alpha := f.WindowColorAlpha
	f.WindowColor = &RGB{
		R: red,
		G: green,
		B: blue,
	}
	color := f.WindowColor
	style := fmt.Sprintf("background-color: rgba(%d, %d, %d, %f);", color.R, color.G, color.B, alpha)
	if runtime.GOOS == "darwin" && f.WindowColorAlpha < 1.0 {
		style = "background-color: rgba(0, 0, 0, 0);"
	}

	borderSizeString := fmt.Sprintf("%d", f.borderSize*2) + "px"
	borderSizeString2 := fmt.Sprintf("%d", f.borderSize*2+f.windowgap) + "px"

	var roundSizeString string
	if runtime.GOOS == "linux" {
		roundSizeString = fmt.Sprintf("%d", f.borderSize*2) + "px"
	} else {
		roundSizeString = "0px"
	}

	if !f.IsBorderless {
		roundSizeString = "0px"
	}

	f.WindowWidget.SetStyleSheet(fmt.Sprintf(`
			#QFramelessWidget {
				border: 0px solid %s; 
				padding-top: %s; padding-right: %s; padding-bottom: %s; padding-left: %s; 
				border-radius: %s;
				%s; 
			}`, color.Hex(), borderSizeString, borderSizeString2, borderSizeString, borderSizeString2, roundSizeString, style),
	)

	if f.IsBorderless {
		f.SetStyleMask()
	}
}

func NewQToolButtonForNotDarwin(parent widgets.QWidget_ITF) *QToolButtonForNotDarwin {
	iconSize := 15
	marginTB := iconSize / 6
	marginLR := 1
	if runtime.GOOS == "linux" {
		iconSize = 18
		marginLR = int(float64(iconSize) / float64(8))
	} else {
		marginLR = int(float64(iconSize) / float64(2.5))
	}

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
	if runtime.GOOS == "darwin" {
		return
	}
	b.IconBtn.SetObjectName(name)
}

func (b *QToolButtonForNotDarwin) Hide() {
	if runtime.GOOS == "darwin" {
		return
	}
	b.Widget.Hide()
}

func (b *QToolButtonForNotDarwin) Show() {
	if runtime.GOOS == "darwin" {
		return
	}
	b.Widget.Show()
}

func (b *QToolButtonForNotDarwin) SetStyle(color *RGB) {
	var backgroundColor string
	if color == nil {
		backgroundColor = "background-color:none;"
	} else {
		hoverColor := color.Brend(b.f.WindowColor, 0.35)
		backgroundColor = fmt.Sprintf("background-color: rgba(%d, %d, %d, %f);", hoverColor.R, hoverColor.G, hoverColor.B, b.f.WindowColorAlpha)
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
	f.IconMinimize.f = f
	f.IconMinimize.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconMinimize.SetObjectName("IconMinimize")
	f.IconMaximize = NewQToolButtonForNotDarwin(nil)
	f.IconMaximize.f = f
	f.IconMaximize.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconMaximize.SetObjectName("IconMaximize")
	f.IconRestore = NewQToolButtonForNotDarwin(nil)
	f.IconRestore.f = f
	f.IconRestore.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconRestore.SetObjectName("IconRestore")
	f.IconClose = NewQToolButtonForNotDarwin(nil)
	f.IconClose.f = f
	f.IconClose.IconBtn.SetFixedSize2(iconSize, iconSize)
	f.IconClose.SetObjectName("IconClose")

	f.SetIconsStyle(nil)
	f.HideButtons()

	dummyWidget := widgets.NewQLabel(nil, 0)
	dummyWidget.SetFixedWidth(20 * 3)

	f.TitleBarLayout.AddWidget(dummyWidget, 0, 0)
	f.TitleBarLayout.AddWidget(f.TitleLabel, 3, 0)
	f.TitleBarLayout.AddWidget(f.TitleBarBtnWidget, 0, 0)

	f.TitleBarBtnLayout.AddWidget(f.IconMinimize.Widget, 0, 0)
	f.TitleBarBtnLayout.AddWidget(f.IconMaximize.Widget, 0, 0)
	f.TitleBarBtnLayout.AddWidget(f.IconRestore.Widget, 0, 0)
	f.TitleBarBtnLayout.AddWidget(f.IconClose.Widget, 0, 0)

	f.TitleBarBtnWidget.SetFixedWidth(30 * 3)
	f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignRight)
	f.TitleBarLayout.SetAlignment(f.TitleLabel, core.Qt__AlignCenter)
}

func (f *QFramelessWindow) ShowButtons() {
	if runtime.GOOS == "darwin" {
		return
	}
	f.IconMinimize.Show()
	f.IconMaximize.Show()
	f.IconRestore.Show()
	f.IconClose.Show()
}

func (f *QFramelessWindow) HideButtons() {
	if runtime.GOOS == "darwin" {
		return
	}
	f.IconMinimize.Hide()
	f.IconMaximize.Hide()
	f.IconRestore.Hide()
	f.IconClose.Hide()
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
	f.TitleBarLayout.SetContentsMargins(5, 5, 0, 5)

	// btnSizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__ToolButton)
	// f.BtnMinimize = widgets.NewQToolButton(f.TitleBarBtnWidget)
	// f.BtnMinimize.SetObjectName("BtnMinimize")
	// f.BtnMinimize.SetSizePolicy(btnSizePolicy)

	// f.BtnMaximize = widgets.NewQToolButton(f.TitleBarBtnWidget)
	// f.BtnMaximize.SetObjectName("BtnMaximize")
	// f.BtnMaximize.SetSizePolicy(btnSizePolicy)

	// f.BtnRestore = widgets.NewQToolButton(f.TitleBarBtnWidget)
	// f.BtnRestore.SetObjectName("BtnRestore")
	// f.BtnRestore.SetSizePolicy(btnSizePolicy)
	// f.BtnRestore.SetVisible(false)

	// f.BtnClose = widgets.NewQToolButton(f.TitleBarBtnWidget)
	// f.BtnClose.SetObjectName("BtnClose")
	// f.BtnClose.SetSizePolicy(btnSizePolicy)

	// NOTE: We use native buttons
	// dummyWidget := widgets.NewQLabel(nil, 0)
	// dummyWidget.SetFixedWidth(20 * 3)

	f.TitleBarLayout.SetSpacing(0)

	// NOTE: We use native buttons
	// f.TitleBarLayout.AddWidget(f.TitleBarBtnWidget, 0, 0)
	// f.TitleBarBtnLayout.AddWidget(f.BtnClose, 0, 0)
	// f.TitleBarBtnLayout.AddWidget(f.BtnMinimize, 0, 0)
	// f.TitleBarBtnLayout.AddWidget(f.BtnMaximize, 0, 0)
	// f.TitleBarBtnLayout.AddWidget(f.BtnRestore, 0, 0)

	f.TitleBarLayout.AddWidget(f.TitleLabel, 2, 0)
	// f.TitleBarLayout.AddWidget(dummyWidget, 0, 0)

	// f.TitleBarBtnWidget.SetFixedWidth(20 * 3)
	// f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignLeft)
	// f.TitleBarLayout.SetAlignment(f.TitleLabel, core.Qt__AlignCenter)
}

func (f *QFramelessWindow) SetupAttributes() {
	if runtime.GOOS == "darwin" {
		f.SetAttribute(core.Qt__WA_TranslucentBackground, true)
		return
	}
	f.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	// f.SetAttribute(core.Qt__WA_NoSystemBackground, true)
	f.SetAttribute(core.Qt__WA_Hover, true)
	f.SetMouseTracking(true)
}

func (f *QFramelessWindow) SetupWindowFlags() {
	defer f.SetStyleMask()
	if runtime.GOOS == "darwin" {
		return
	}
	// if runtime.GOOS == "windows" {
	// 	f.SetWindowFlag(core.Qt__Window, true)
	// } else {
	// 	f.SetWindowFlag(core.Qt__Window, false)
	// }

	f.SetWindowFlag(core.Qt__FramelessWindowHint, true)
	// f.SetWindowFlag(core.Qt__NoDropShadowWindowHint, true)
	// f.SetWindowFlag(core.Qt__CustomizeWindowHint, true)
	// f.SetWindowFlag(core.Qt__WindowTitleHint, false)

	// if runtime.GOOS == "windows" {
	// 	f.SetWindowFlag(core.Qt__WindowMaximizeButtonHint, true)
	// }
}

func (f *QFramelessWindow) SetupTitleIcon(filename string) {
	// s := 14
	// f.TitleIconLabel.SetMaximumWidth(s)
	// f.TitleIconLabel.SetMinimumWidth(s)
	// pic := gui.NewQPixmap3(filename, "", core.Qt__AutoColor)
	// pic = pic.Scaled2(s, s, core.Qt__KeepAspectRatio, core.Qt__SmoothTransformation)
	// f.TitleIconLabel.SetPixmap(pic)
	// f.TitleLabel.SetFixedWidth(f.TitleStringLabel.FontMetrics().BoundingRect2(f.TitleStringLabel.Text()).Width()+f.TitleIconLabel.Width())
	// f.TitleIconLabel.Show()
	// f.IsTitleIconShown = true
}

func (f *QFramelessWindow) SetupTitle(title string) {
	if !f.IsBorderless {
		f.SetWindowTitle(title)
		return
	}
	f.TitleStringLabel.SetText(title)
	f.TitleStringLabel.SetFixedWidth(f.TitleStringLabel.FontMetrics().BoundingRect2(title).Width())
	f.TitleStringLabel.SetSizePolicy2(widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Minimum)
	if f.IsTitleIconShown {
		f.TitleLabel.SetFixedWidth(f.TitleStringLabel.FontMetrics().BoundingRect2(title).Width() + f.TitleIconLabel.Width())
	} else {
		f.TitleLabel.SetFixedWidth(f.TitleStringLabel.FontMetrics().BoundingRect2(title).Width())
	}
}

func (f *QFramelessWindow) SetupTitleColor(red uint16, green uint16, blue uint16) {
	if !f.IsBorderless {
		return
	}
	f.TitleColor = &RGB{
		R: red,
		G: green,
		B: blue,
	}
	f.SetupTitleBarColor()
}

func (f *QFramelessWindow) SetupTitleBarColor() {
	if !f.IsBorderless {
		return
	}
	var color, labelColor *RGB
	brendRatio := 0.0
	if f.IsActiveWindow() {
		color = f.TitleColor
	} else {
		brendRatio = 0.15
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
		f.TitleLabel.SetStyleSheet(fmt.Sprintf(" * { color: rgb(%d, %d, %d); }", labelColor.R, labelColor.G, labelColor.B))
		f.SetupTitleBarColorForNotDarwin(
			&RGB{
				R: uint16((float64(f.WindowColor.R)*brendRatio + float64(labelColor.R)*(1.0-brendRatio))),
				G: uint16((float64(f.WindowColor.G)*brendRatio + float64(labelColor.G)*(1.0-brendRatio))),
				B: uint16((float64(f.WindowColor.B)*brendRatio + float64(labelColor.B)*(1.0-brendRatio))),
			},
		)
	} else {
		f.TitleLabel.SetStyleSheet(fmt.Sprintf(" * { color: rgb(%d, %d, %d); }", labelColor.R, labelColor.G, labelColor.B))
	}

	if runtime.GOOS == "linux" {
		f.SetupTitleBarActions()
	}
}

func (f *QFramelessWindow) SetupTitleBarColorForNotDarwin(color *RGB) {
	var SvgMinimize, SvgMaximize, SvgRestore, SvgClose string
	var closeColor *RGB

	if runtime.GOOS == "windows" {
		SvgMinimize = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M20,14H4V10H20" />
		</svg>
		`, color.Hex())

		SvgMaximize = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M4,4H20V20H4V4M6,8V18H18V8H6Z" />
		</svg>
		`, color.Hex())

		SvgRestore = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M4,8H8V4H20V16H16V20H4V8M16,8V14H18V6H10V8H16M6,12V18H14V12H6Z" />
		</svg>
		`, color.Hex())

		SvgClose = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M13.46,12L19,17.54V19H17.54L12,13.46L6.46,19H5V17.54L10.54,12L5,6.46V5H6.46L12,10.54L17.54,5H19V6.46L13.46,12Z" />
		</svg>
		`, color.Hex())
	} else {
		if f.IsActiveWindow() {
			closeColor = &RGB{
				R: 232,
				G: 96,
				B: 50,
			}
		} else {
			closeColor = color
		}

		SvgMinimize = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M17,13H7V11H17M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
		</svg>
		`, color.Hex())

		SvgMaximize = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
		<g transform="scale(0.6) translate(8,8)">
			<path fill="%s" d="M19,3H5C3.89,3 3,3.89 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19V5C21,3.89 20.1,3 19,3M19,5V19H5V5H19Z" />
		</g>
		</svg>
		`, color.Hex(), f.WindowColor.Hex())

		SvgRestore = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
		<g transform="scale(0.6) translate(8,8)">
			<path fill="%s" d="M19,3H5C3.89,3 3,3.89 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19V5C21,3.89 20.1,3 19,3M19,5V19H5V5H19Z" />
		</g>
		</svg>
		`, color.Hex(), f.WindowColor.Hex())

		SvgClose = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<g transform="translate(0,1)">
		<path fill="%s" d="M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"/><path d="M0 0h24v24H0z" fill="none"/></g></svg>
		`, closeColor.Hex())
	}

	f.IconMinimize.IconBtn.Load2(core.NewQByteArray2(SvgMinimize, len(SvgMinimize)))
	f.IconMaximize.IconBtn.Load2(core.NewQByteArray2(SvgMaximize, len(SvgMaximize)))
	f.IconRestore.IconBtn.Load2(core.NewQByteArray2(SvgRestore, len(SvgRestore)))
	f.IconClose.IconBtn.Load2(core.NewQByteArray2(SvgClose, len(SvgClose)))

	f.IconRestore.Widget.SetVisible(false)
}

func (f *QFramelessWindow) SetupContent(layout widgets.QLayout_ITF) {
	if f.Content != nil {
		f.Content.SetLayout(layout)
	} else {
		f.WindowWidget.SetLayout(layout)
	}
}

func (f *QFramelessWindow) UpdateWidget() {
	f.Widget.Update()
	f.Update()
}

func (f *QFramelessWindow) SetupWindowActions() {
	// Ref: https://stackoverflow.com/questions/5752408/qt-resize-borderless-widget/37507341#37507341
	f.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
		return f.QFramelessDefaultEventFilter(watched, event)
	})
}

func (f *QFramelessWindow) QFramelessDefaultEventFilter(watched *core.QObject, event *core.QEvent) bool {
	if f.IsBorderless {
		switch event.Type() {
		case core.QEvent__ActivationChange:
			f.SetupTitleBarColor()

		case core.QEvent__Resize:
			if runtime.GOOS == "darwin" {
				f.SetStyleMask()
			}
		case core.QEvent__WindowStateChange:
			if runtime.GOOS == "darwin" {
				e := gui.NewQWindowStateChangeEventFromPointer(core.PointerFromQEvent(event))
				if e.OldState() == core.Qt__WindowFullScreen {
					f.SetStyleMask()
					f.ShowTitlebar()
					f.ActivateWindow()
				}
				if f.WindowState() == core.Qt__WindowFullScreen {
					f.SetStyleMask()
					f.HideTitlebar()
				}
			}
			if runtime.GOOS == "windows" {
				// It is a workaround for https://github.com/akiyosi/goneovim/issues/91#issuecomment-587041657
				// e := gui.NewQWindowStateChangeEventFromPointer(core.PointerFromQEvent(event))
				// if e.OldState() == core.Qt__WindowMinimized {
				if f.WindowState() == core.Qt__WindowMinimized {
					go func() {
						time.Sleep(300 * time.Millisecond)

						frameRect := f.WindowWidget.FrameGeometry()

						left := 0
						top := 0
						right := frameRect.Width()
						bottom := frameRect.Height()
						topLeftPoint := core.NewQPoint2(left, top)
						rightBottomPoint := core.NewQPoint2(right, bottom)
						f.WindowWidget.SetGeometry(
							core.NewQRect2(
								topLeftPoint,
								rightBottomPoint,
							),
						)

						frameRect = f.WindowWidget.FrameGeometry()
						rect := f.FrameGeometry()
						frameWidth := frameRect.Width() - rect.Width()
						frameHeight := frameRect.Height() - rect.Height()

						left = rect.Left() - frameWidth/2
						top = rect.Top()
						right = rect.Right() + frameWidth/2
						bottom = rect.Bottom() + frameHeight
						topLeftPoint = core.NewQPoint2(left, top)
						rightBottomPoint = core.NewQPoint2(right, bottom)
						f.SetGeometry(
							core.NewQRect2(
								topLeftPoint,
								rightBottomPoint,
							),
						)
					}()
				}
			}

		case core.QEvent__HoverMove:
			e := gui.NewQMouseEventFromPointer(core.PointerFromQEvent(event))
			if runtime.GOOS == "darwin" {
				f.showButtonsInDarwin(e.GlobalPos(), f.FrameGeometry())
			}
			if runtime.GOOS == "windows" {
				f.showTitlebarInWindows(e.GlobalPos(), f.FrameGeometry())
			}

			f.updateCursorShape(e.GlobalPos())

		case core.QEvent__Leave:
			if runtime.GOOS == "darwin" {
				return f.Widget.EventFilter(watched, event)
			}

			cursor := gui.NewQCursor()
			cursor.SetShape(core.Qt__ArrowCursor)
			f.SetCursor(cursor)

		case core.QEvent__MouseMove:
			if runtime.GOOS == "darwin" {
				return f.Widget.EventFilter(watched, event)
			}

			e := gui.NewQMouseEventFromPointer(core.PointerFromQEvent(event))
			f.mouseMove(e)

		case core.QEvent__MouseButtonPress:
			if runtime.GOOS == "darwin" {
				return f.Widget.EventFilter(watched, event)
			}

			e := gui.NewQMouseEventFromPointer(core.PointerFromQEvent(event))
			f.mouseButtonPressed(e)

		case core.QEvent__MouseButtonRelease:
			if runtime.GOOS == "darwin" {
				return f.Widget.EventFilter(watched, event)
			}

			f.isDragStart = false
			f.isLeftButtonPressed = false
			f.hoverEdge = None

		default:
		}
	}

	return f.Widget.EventFilter(watched, event)
}

func (f *QFramelessWindow) showTitlebarInWindows(pos *core.QPoint, rect *core.QRect) {
	if !f.IsTitlebarHidden {
		return
	}

	rectShowingButtons := core.NewQRect4(
		rect.TopLeft().X(),
		rect.TopLeft().Y()+5,
		rect.Width(),
		30,
	)
	if rectShowingButtons.Contains(pos, true) {
		f.TitleBar.Show()
	} else {
		f.TitleBar.Hide()
	}
}

func (f *QFramelessWindow) showButtonsInDarwin(pos *core.QPoint, rect *core.QRect) {
	if !f.IsTitlebarHidden {
		return
	}
	if f.WindowState() == core.Qt__WindowFullScreen {
		f.SetNSWindowStyleMask(
			true,
			f.WindowColor.R, f.WindowColor.G, f.WindowColor.B,
			float32(f.WindowColorAlpha),
			f.WindowState() == core.Qt__WindowFullScreen,
		)

		return
	}

	rectShowingButtons := core.NewQRect4(
		rect.TopLeft().X(),
		rect.TopLeft().Y()+5,
		rect.Width(),
		30,
	)
	if rectShowingButtons.Contains(pos, true) {
		f.SetNSWindowStyleMask(
			true,
			f.WindowColor.R, f.WindowColor.G, f.WindowColor.B,
			float32(f.WindowColorAlpha),
			f.WindowState() == core.Qt__WindowFullScreen,
		)
	} else {
		f.SetNSWindowStyleMask(
			false,
			f.WindowColor.R, f.WindowColor.G, f.WindowColor.B,
			float32(f.WindowColorAlpha),
			f.WindowState() == core.Qt__WindowFullScreen,
		)
	}
}

func (f *QFramelessWindow) mouseMove(e *gui.QMouseEvent) {
	// https://stackoverflow.com/questions/5752408/qt-resize-borderless-widget/37507341

	if f.isLeftButtonPressed {

		if f.hoverEdge != None {
			X := e.GlobalPos().X()
			Y := e.GlobalPos().Y()

			if f.MousePos[0] == X && f.MousePos[1] == Y {
				return
			}

			f.MousePos[0] = X
			f.MousePos[1] = Y

			var left, top, right, bottom int
			var topLeftPoint, rightBottomPoint *core.QPoint

			if runtime.GOOS == "windows" {
				// Use tricky workaround only on Windows due to the following issues
				// https://github.com/therecipe/qt/issues/938
				frameRect := f.WindowWidget.FrameGeometry()
				rect := f.FrameGeometry()
				frameWidth := frameRect.Width() - rect.Width()
				frameHeight := frameRect.Height() - rect.Height()
				left = rect.Left() - frameWidth/2
				top = rect.Top()
				right = rect.Right() + frameWidth/2
				bottom = rect.Bottom() + frameHeight
			} else {
				rect := f.FrameGeometry()
				left = rect.Left()
				top = rect.Top()
				right = rect.Right()
				bottom = rect.Bottom()
			}

			switch f.hoverEdge {
			case Top:
				top = Y
			case Bottom:
				bottom = Y
			case Left:
				left = X
			case Right:
				right = X
			case TopLeft:
				top = Y
				left = X
			case TopRight:
				top = Y
				right = X
			case BottomLeft:
				bottom = Y
				left = X
			case BottomRight:
				bottom = Y
				right = X
			default:
			}

			topLeftPoint = core.NewQPoint2(left, top)
			rightBottomPoint = core.NewQPoint2(right, bottom)
			rect := core.NewQRect2(topLeftPoint, rightBottomPoint)

			// minimum size
			minimumWidth := f.minimumWidth
			minimumHeight := f.minimumHeight
			if rect.Width() <= minimumWidth {
				switch f.hoverEdge {
				case Left:
					left = right - minimumWidth
				case Right:
					right = left + minimumWidth
				case TopLeft:
					left = right - minimumWidth
				case TopRight:
					right = left + minimumWidth
				case BottomLeft:
					left = right - minimumWidth
				case BottomRight:
					right = left + minimumWidth
				default:
				}
			}
			if rect.Height() <= minimumHeight {
				switch f.hoverEdge {
				case Top:
					top = bottom - minimumHeight
				case Bottom:
					bottom = top + minimumHeight
				case TopLeft:
					top = bottom - minimumHeight
				case TopRight:
					top = bottom - minimumHeight
				case BottomLeft:
					bottom = top + minimumHeight
				case BottomRight:
					bottom = top + minimumHeight
				default:
				}
			}
			if rect.Width() <= minimumWidth || rect.Height() <= minimumHeight {
				right = left + minimumWidth
				bottom = top + minimumHeight
			}

			topLeftPoint = core.NewQPoint2(left, top)
			rightBottomPoint = core.NewQPoint2(right, bottom)
			newRect := core.NewQRect2(topLeftPoint, rightBottomPoint)

			f.SetGeometry(newRect)
		}
	}
}

func (f *QFramelessWindow) mouseButtonPressed(e *gui.QMouseEvent) {
	f.hoverEdge = f.calcCursorPos(e.GlobalPos(), f.FrameGeometry())
	if f.hoverEdge != None {
		f.isLeftButtonPressed = true
	}
}

func (f *QFramelessWindow) mouseButtonPressedForWin(e *gui.QMouseEvent) {
	if f.hoverEdge != None {
		f.isLeftButtonPressed = true
	}
}

func (f *QFramelessWindow) updateCursorShape(pos *core.QPoint) {
	if f.isLeftButtonPressed {
		return
	}
	window := f
	cursor := gui.NewQCursor()
	if window.IsFullScreen() || window.IsMaximized() {
		if f.isCursorChanged {
			cursor.SetShape(core.Qt__ArrowCursor)
			window.SetCursor(cursor)
		}
	}
	hoverEdge := f.calcCursorPos(pos, window.FrameGeometry())
	f.isCursorChanged = true
	switch hoverEdge {
	case Top, Bottom:
		cursor.SetShape(core.Qt__SizeVerCursor)
		window.SetCursor(cursor)
	case Left, Right:
		cursor.SetShape(core.Qt__SizeHorCursor)
		window.SetCursor(cursor)
	case TopLeft, BottomRight:
		cursor.SetShape(core.Qt__SizeFDiagCursor)
		window.SetCursor(cursor)
	case TopRight, BottomLeft:
		cursor.SetShape(core.Qt__SizeBDiagCursor)
		window.SetCursor(cursor)
	default:
		cursor.SetShape(core.Qt__ArrowCursor)
		window.SetCursor(cursor)
		f.isCursorChanged = false
	}
}

func (f *QFramelessWindow) calcCursorPos(pos *core.QPoint, rect *core.QRect) Edge {
	frameWidth := f.WindowWidget.Width() - rect.Width()
	frameHeight := f.WindowWidget.Height() - rect.Height()
	rectX := rect.X() - frameWidth/2 + 1
	rectY := rect.Y() - frameHeight/2 + 1
	rectWidth := rect.Width() + frameWidth - 1
	rectHeight := rect.Height() + frameHeight - 1
	posX := pos.X()
	posY := pos.Y()

	edge := f.detectEdgeOnCursor(posX, posY, rectX, rectY, rectWidth, rectHeight)
	return edge
}

func (f *QFramelessWindow) detectEdgeOnCursor(posX, posY, rectX, rectY, rectWidth, rectHeight int) Edge {
	doubleBorderSize := f.borderSize * 2
	octupleBorderSize := f.borderSize * 8
	topBorderSize := 2 - 1

	var onLeft, onRight, onBottom, onTop, onBottomLeft, onBottomRight, onTopRight, onTopLeft bool

	onBottomLeft = (((posX <= (rectX + octupleBorderSize)) && posX >= rectX &&
		(posY <= (rectY + rectHeight)) && (posY >= (rectY + rectHeight - doubleBorderSize))) ||
		((posX <= (rectX + doubleBorderSize)) && posX >= rectX &&
			(posY <= (rectY + rectHeight)) && (posY >= (rectY + rectHeight - octupleBorderSize))))

	if onBottomLeft {
		return BottomLeft
	}

	onBottomRight = (((posX >= (rectX + rectWidth - octupleBorderSize)) && (posX <= (rectX + rectWidth)) &&
		(posY >= (rectY + rectHeight - doubleBorderSize)) && (posY <= (rectY + rectHeight))) ||
		((posX >= (rectX + rectWidth - doubleBorderSize)) && (posX <= (rectX + rectWidth)) &&
			(posY >= (rectY + rectHeight - octupleBorderSize)) && (posY <= (rectY + rectHeight))))

	if onBottomRight {
		return BottomRight
	}

	onTopRight = (posX >= (rectX + rectWidth - doubleBorderSize)) && (posX <= (rectX + rectWidth)) &&
		(posY >= rectY) && (posY <= (rectY + doubleBorderSize))
	if onTopRight {
		return TopRight
	}

	onTopLeft = posX >= rectX && (posX <= (rectX + doubleBorderSize)) &&
		posY >= rectY && (posY <= (rectY + doubleBorderSize))
	if onTopLeft {
		return TopLeft
	}

	onLeft = (posX >= (rectX - doubleBorderSize)) && (posX <= (rectX + doubleBorderSize)) &&
		(posY <= (rectY + rectHeight - doubleBorderSize)) &&
		(posY >= rectY+doubleBorderSize)
	if onLeft {
		return Left
	}

	onRight = (posX >= (rectX + rectWidth - doubleBorderSize)) &&
		(posX <= (rectX + rectWidth)) &&
		(posY >= (rectY + doubleBorderSize)) && (posY <= (rectY + rectHeight - doubleBorderSize))
	if onRight {
		return Right
	}

	onBottom = (posX >= (rectX + doubleBorderSize)) && (posX <= (rectX + rectWidth - doubleBorderSize)) &&
		(posY >= (rectY + rectHeight - doubleBorderSize)) && (posY <= (rectY + rectHeight))
	if onBottom {
		return Bottom
	}

	onTop = (posX >= (rectX + doubleBorderSize)) && (posX <= (rectX + rectWidth - doubleBorderSize)) &&
		(posY >= rectY) && (posY <= (rectY + topBorderSize))
	if onTop {
		return Top
	}

	return None
}

func (f *QFramelessWindow) HideTitlebar() {
	f.TitleBar.Hide()
}

func (f *QFramelessWindow) ShowTitlebar() {
	if f.IsTitlebarHidden {
		f.TitleBar.Hide()
		return
	}
	f.TitleBar.Show()
}

func (c *RGB) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", uint8(c.R), uint8(c.G), uint8(c.B))
}

func (c *RGB) Brend(color *RGB, alpha float64) *RGB {
	if color == nil {
		return &RGB{0, 0, 0}
	}
	return &RGB{
		R: uint16((float64(c.R) * float64(1-alpha)) + (float64(color.R) * float64(alpha))),
		G: uint16((float64(c.G) * float64(1-alpha)) + (float64(color.G) * float64(alpha))),
		B: uint16((float64(c.B) * float64(1-alpha)) + (float64(color.B) * float64(alpha))),
	}
}

func (c *RGB) toColorref() uint32 {
	b := uint32(c.B) << 16
	g := uint32(c.G) << 8
	r := uint32(c.R)

	return b + g + r
}
