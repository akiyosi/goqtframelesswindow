package qframelesswindow

import (
	"fmt"
	"time"
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
	f       *QFramelessWindow
	Widget  *widgets.QWidget
	IconBtn *svg.QSvgWidget
	isHover bool
}

type QFramelessWindow struct {
	widgets.QMainWindow

	WindowColor      *RGB
	WindowColorAlpha float64

	Widget  *widgets.QWidget
	Layout  *widgets.QVBoxLayout
	Content *widgets.QWidget

	WindowWidget  *widgets.QFrame
	WindowVLayout *widgets.QVBoxLayout
	shadowMargin  int
	borderSize    int
	minimumWidth  int
	minimumHeight int

	TitleBar          *widgets.QWidget
	TitleBarLayout    *widgets.QHBoxLayout
	// TitleLabel        *widgets.QLabel
	TitleIconLabel    *widgets.QLabel
	TitleStringLabel  *widgets.QLabel
	TitleLabel        *widgets.QWidget
	TitleBarBtnWidget *widgets.QWidget
	TitleBarBtnLayout *widgets.QHBoxLayout
	TitleColor        *RGB
	TitleBarMousePos  *core.QPoint
	IsTitleBarPressed bool
	IsTitleIconShown  bool

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
	hoverEdge           Edge
	Position            *core.QPoint
	MousePos            [2]int

	borderless bool
}

func CreateQFramelessWindow(alpha float64) *QFramelessWindow {
	f := NewQFramelessWindow(nil, 0)
	f.WindowColorAlpha = alpha
	if f.WindowColorAlpha == 1.0 {
		f.SetupNativeEvent()
	} else {
		f.SetupNativeEvent2()
	}
	f.Widget = widgets.NewQWidget(nil, 0)
	f.Widget.SetStyleSheet(" * { background-color: rgba(0, 0, 0, 0.0); color: rgba(0, 0, 0, 0); }")
	f.SetCentralWidget(f.Widget)

	f.shadowMargin = 0
	f.SetupBorderSize(3)
	f.SetupUI(f.Widget)
	f.SetupWindowFlags()
	f.SetupAttributes()
	f.SetupWindowActions()
	f.SetupTitleBarActions()
	f.SetupMinimumSize(400, 300)

	return f
}

func (f *QFramelessWindow) SetupMinimumSize(w int, h int) {
	W := w + (2 * f.shadowMargin)
	H := h + (2 * f.shadowMargin)
	f.SetMinimumSize2(W, H)
	f.Widget.SetMinimumSize2(W, H)
	f.WindowWidget.SetMinimumSize2(w, h)
	f.minimumWidth = w
	f.minimumHeight = h
}

func (f *QFramelessWindow) SetupBorderSize(size int) {
	f.borderSize = size
}

// For MacOS only
func (f *QFramelessWindow) AddWindowNativeShadow() {
	f.SetWindowFlag(core.Qt__NoDropShadowWindowHint, false)
	f.SetStyleMask()
}

// For MacOS only
func (f *QFramelessWindow) RemoveWindowNativeShadow() {
	f.SetWindowFlag(core.Qt__NoDropShadowWindowHint, true)
	f.SetStyleMask()
}

func (f *QFramelessWindow) SetupUI(widget *widgets.QWidget) {
	f.InstallEventFilter(f)

	widget.SetSizePolicy2(widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum, widgets.QSizePolicy__Expanding|widgets.QSizePolicy__Maximum)
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

	// create window content
	f.Content = widgets.NewQWidget(f.WindowWidget, 0)

	// Set widget to layout
	f.WindowVLayout.AddWidget(f.TitleBar, 0, 0)
	f.WindowVLayout.AddWidget(f.Content, 0, 0)

	f.Layout.AddWidget(f.WindowWidget, 0, 0)
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
	borderSizeString := fmt.Sprintf("%d", f.borderSize*2) + "px"

	var roundSizeString string
	if runtime.GOOS != "windows" {
		roundSizeString = fmt.Sprintf("%d", f.borderSize*2) + "px"
	} else {
		roundSizeString = "0px"
	}

	f.WindowWidget.SetStyleSheet(fmt.Sprintf(`
	#QFramelessWidget {
		border: 0px solid %s; 
		padding-top: 2px; padding-right: %s; padding-bottom: %s; padding-left: %s; 
		border-radius: %s;
		%s; 
	}`, color.Hex(), borderSizeString, borderSizeString, borderSizeString, roundSizeString, style))
	f.SetStyleMask()
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
		hoverColor := color.Brend(b.f.WindowColor, 0.65)
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

	f.IconMinimize.Hide()
	f.IconMaximize.Hide()
	f.IconRestore.Hide()
	f.IconClose.Hide()

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
	btnSizePolicy := widgets.NewQSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__ToolButton)
	f.BtnMinimize = widgets.NewQToolButton(f.TitleBarBtnWidget)
	f.BtnMinimize.SetObjectName("BtnMinimize")
	f.BtnMinimize.SetSizePolicy(btnSizePolicy)

	f.BtnMaximize = widgets.NewQToolButton(f.TitleBarBtnWidget)
	f.BtnMaximize.SetObjectName("BtnMaximize")
	f.BtnMaximize.SetSizePolicy(btnSizePolicy)

	f.BtnRestore = widgets.NewQToolButton(f.TitleBarBtnWidget)
	f.BtnRestore.SetObjectName("BtnRestore")
	f.BtnRestore.SetSizePolicy(btnSizePolicy)
	f.BtnRestore.SetVisible(false)

	f.BtnClose = widgets.NewQToolButton(f.TitleBarBtnWidget)
	f.BtnClose.SetObjectName("BtnClose")
	f.BtnClose.SetSizePolicy(btnSizePolicy)

	dummyWidget := widgets.NewQLabel(nil, 0)
	dummyWidget.SetFixedWidth(20 * 3)

	f.TitleBarLayout.SetSpacing(0)
	f.TitleBarLayout.AddWidget(f.TitleBarBtnWidget, 0, 0)
	f.TitleBarLayout.AddWidget(f.TitleLabel, 3, 0)
	f.TitleBarLayout.AddWidget(dummyWidget, 0, 0)

	f.TitleBarBtnLayout.AddWidget(f.BtnClose, 0, 0)
	f.TitleBarBtnLayout.AddWidget(f.BtnMinimize, 0, 0)
	f.TitleBarBtnLayout.AddWidget(f.BtnMaximize, 0, 0)
	f.TitleBarBtnLayout.AddWidget(f.BtnRestore, 0, 0)

	f.TitleBarBtnWidget.SetFixedWidth(20 * 3)
	f.TitleBarLayout.SetAlignment(f.TitleBarBtnWidget, core.Qt__AlignLeft)
	f.TitleBarLayout.SetAlignment(f.TitleLabel, core.Qt__AlignCenter)
}

func (f *QFramelessWindow) SetupAttributes() {
	f.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	f.SetAttribute(core.Qt__WA_NoSystemBackground, true)
	f.SetAttribute(core.Qt__WA_Hover, true)
	f.SetMouseTracking(true)
}

func (f *QFramelessWindow) SetupWindowFlags() {
	f.SetWindowFlag(core.Qt__Window, true)
	f.SetWindowFlag(core.Qt__FramelessWindowHint, true)
	f.SetWindowFlag(core.Qt__NoDropShadowWindowHint, true)
	if runtime.GOOS == "linux" {
		f.SetWindowFlag(core.Qt__Window, false)
		f.SetWindowFlag(core.Qt__WindowStaysOnTopHint, true)
	}
	f.SetStyleMask()
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
	f.TitleStringLabel.SetText(title)
	f.TitleStringLabel.SetFixedWidth(f.TitleStringLabel.FontMetrics().BoundingRect2(title).Width())
	f.TitleStringLabel.SetSizePolicy2(widgets.QSizePolicy__Minimum, widgets.QSizePolicy__Minimum)
	if f.IsTitleIconShown {
		f.TitleLabel.SetFixedWidth(f.TitleStringLabel.FontMetrics().BoundingRect2(title).Width()+f.TitleIconLabel.Width())
	} else {
		f.TitleLabel.SetFixedWidth(f.TitleStringLabel.FontMetrics().BoundingRect2(title).Width())
	}
}

func (f *QFramelessWindow) SetupTitleColor(red uint16, green uint16, blue uint16) {
	f.TitleColor = &RGB{
		R: red,
		G: green,
		B: blue,
	}
	f.SetupTitleBarColor()
}

func (f *QFramelessWindow) SetupTitleBarColor() {
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
				R: uint16((float64(f.WindowColor.R) * brendRatio + float64(labelColor.R) * (1.0 - brendRatio))),
				G: uint16((float64(f.WindowColor.G) * brendRatio + float64(labelColor.G) * (1.0 - brendRatio))),
				B: uint16((float64(f.WindowColor.B) * brendRatio + float64(labelColor.B) * (1.0 - brendRatio))),
			},
		)
	} else {
		f.TitleLabel.SetStyleSheet(fmt.Sprintf(" * { color: rgb(%d, %d, %d); }", labelColor.R, labelColor.G, labelColor.B))
		f.SetupTitleBarColorForDarwin(color)
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

	f.IconMinimize.Show()
	f.IconMaximize.Show()
	f.IconRestore.Show()
	f.IconRestore.Widget.SetVisible(false)
	f.IconClose.Show()
}

func (f *QFramelessWindow) SetupTitleBarColorForDarwin(color *RGB) {
	var baseStyle, restoreAndMaximizeColor, minimizeColor, closeColor string
	baseStyle = ` #BtnMinimize, #BtnMaximize, #BtnRestore, #BtnClose {
		min-width: 12px;
		min-height: 12px;
		max-width: 12px;
		max-height: 12px;
		border-radius: 6px;
		border-width: 0px;
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
				border-color: rgba(128, 128, 128, 0.2);
			}
		`
		minimizeColor = `
			#BtnMinimize {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgba(128, 128, 128, 0.2);
			}
		`
		closeColor = `
			#BtnClose {
				background-color: rgba(128, 128, 128, 0.3);
				border-color: rgba(128, 128, 128, 0.2);
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

func (f *QFramelessWindow) SetupContent(layout widgets.QLayout_ITF) {
	f.Content.SetLayout(layout)
}

func (f *QFramelessWindow) UpdateWidget() {
	f.Widget.Update()
	f.Update()
}

func (f *QFramelessWindow) SetupWindowActions() {
	// Ref: https://stackoverflow.com/questions/5752408/qt-resize-borderless-widget/37507341#37507341
	f.ConnectEventFilter(func(watched *core.QObject, event *core.QEvent) bool {
		e := gui.NewQMouseEventFromPointer(core.PointerFromQEvent(event))
		switch event.Type() {
		case core.QEvent__ActivationChange:
			f.SetupTitleBarColor()

		case core.QEvent__WindowStateChange:
			rect := f.WindowWidget.FrameGeometry()
			// It is a workaround for https://github.com/akiyosi/goneovim/issues/91#issuecomment-587041657
			if f.WindowState() == core.Qt__WindowMinimized {
				if !(rect.Width() == f.minimumWidth && rect.Height() == f.minimumHeight) {
					go func() {
						time.Sleep(300 * time.Millisecond)
						f.SetGeometry(f.FrameGeometry())
					}()
				}
			}

		case core.QEvent__HoverMove:
			f.updateCursorShape(e.GlobalPos())

		case core.QEvent__Leave:
			cursor := gui.NewQCursor()
			cursor.SetShape(core.Qt__ArrowCursor)
			f.SetCursor(cursor)

		case core.QEvent__MouseMove:
			f.mouseMove(e)

		case core.QEvent__MouseButtonPress:
			f.mouseButtonPressed(e)

		case core.QEvent__MouseButtonRelease:
			f.isDragStart = false
			f.isLeftButtonPressed = false
			f.hoverEdge = None

		default:
		}

		return f.Widget.EventFilter(watched, event)
	})
}

func (f *QFramelessWindow) mouseMove(e *gui.QMouseEvent) {
	// https://stackoverflow.com/questions/5752408/qt-resize-borderless-widget/37507341
	window := f
	// margin := f.shadowMargin

	if f.isLeftButtonPressed {

		if f.hoverEdge != None {
			X := e.GlobalPos().X()
			Y := e.GlobalPos().Y()

			if f.MousePos[0] == X && f.MousePos[1] == Y {
				return
			}

			f.MousePos[0] = X
			f.MousePos[1] = Y

			var left, top, right, bottom, frameWidth, frameHeight int
			var topLeftPoint, rightBottomPoint *core.QPoint
			if runtime.GOOS == "windows" {
				// Use tricky workaround only on Windows due to the following issues
				// https://github.com/therecipe/qt/issues/938
				entireRect := f.WindowWidget.FrameGeometry()
				innerRect := f.FrameGeometry()
				frameWidth = entireRect.Width() - innerRect.Width()
				frameHeight = entireRect.Height() - innerRect.Height()
				left =   innerRect.Left() - frameWidth/2
				top =    innerRect.Top()
				right =  innerRect.Right() + frameWidth/2
				bottom = innerRect.Bottom() + frameHeight
			} else {
				entireRect := window.FrameGeometry()
				left =   entireRect.Left()
				top =    entireRect.Top()
				right =  entireRect.Right()
				bottom = entireRect.Bottom()
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

			window.SetGeometry(newRect)
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
	frameWidth :=  f.WindowWidget.Width() - rect.Width()
	frameHeight := f.WindowWidget.Height() - rect.Height()
	rectX := rect.X() - frameWidth/2 + 1
	rectY := rect.Y() - frameHeight/2 + 1
	rectWidth := rect.Width()   + frameWidth - 1
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

	// margin := f.shadowMargin
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
