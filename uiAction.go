// +build !darwin

package qframelesswindow

import (
	"fmt"
	"runtime"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

func (f *QFramelessWindow) SetupTitleBarActions() {
	t := f.TitleBar

	if runtime.GOOS == "windows" {
		// hover style for windows
		f.IconMinimize.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconMinimize.SetStyle(&RGB{
				R: 0,
				G: 162,
				B: 232,
			})
		})
		f.IconMaximize.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconMaximize.SetStyle(&RGB{
				R: 0,
				G: 162,
				B: 232,
			})
		})
		f.IconRestore.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconRestore.SetStyle(&RGB{
				R: 0,
				G: 162,
				B: 232,
			})
		})
		f.IconClose.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconClose.SetStyle(&RGB{
				R: 232,
				G: 0,
				B: 72,
			})
		})

		f.IconMinimize.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconMinimize.SetStyle(nil)
		})
		f.IconMaximize.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconMaximize.SetStyle(nil)
		})
		f.IconRestore.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconRestore.SetStyle(nil)
		})
		f.IconClose.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconClose.SetStyle(nil)
		})
	} else {  // hover style for linux
		brendRatio := 0.65
		hoverBrendRatio := 0.5
		isActive := f.IsActiveWindow()
		if !isActive {
			brendRatio = 0.8
		}
		labelColor := f.TitleColor 
		if labelColor == nil {
			labelColor = &RGB{
				R: 128,
				G: 128,
				B: 128,
			}
		}
		backgroundColor := f.WindowColor
		if backgroundColor == nil {
			backgroundColor = &RGB{
				R: 0,
				G: 0,
				B: 0,
			}

		}
		hoverColor := &RGB{
				R: uint16((float64(backgroundColor.R) * hoverBrendRatio + float64(labelColor.R) * (1.0 - hoverBrendRatio))),
				G: uint16((float64(backgroundColor.G) * hoverBrendRatio + float64(labelColor.G) * (1.0 - hoverBrendRatio))),
				B: uint16((float64(backgroundColor.B) * hoverBrendRatio + float64(labelColor.B) * (1.0 - hoverBrendRatio))),
		}
		color := &RGB{
				R: uint16((float64(backgroundColor.R) * brendRatio + float64(labelColor.R) * (1.0 - brendRatio))),
				G: uint16((float64(backgroundColor.G) * brendRatio + float64(labelColor.G) * (1.0 - brendRatio))),
				B: uint16((float64(backgroundColor.B) * brendRatio + float64(labelColor.B) * (1.0 - brendRatio))),
		}
		hoverCloseColor := &RGB{
			R: 255,
			G: 111,
			B: 65,
		}
		closeColor := &RGB{
			R: 232,
			G: 96,
			B: 50,
		}
		if !isActive {
			closeColor = color
		}


		f.IconMinimize.Widget.DisconnectEnterEvent()
		f.IconMaximize.Widget.DisconnectEnterEvent()
		f.IconRestore.Widget.DisconnectEnterEvent()
		f.IconClose.Widget.DisconnectEnterEvent()
		f.IconMinimize.Widget.DisconnectLeaveEvent()
		f.IconMaximize.Widget.DisconnectLeaveEvent()
		f.IconRestore.Widget.DisconnectLeaveEvent()
		f.IconClose.Widget.DisconnectLeaveEvent()
		f.IconMinimize.Widget.DisconnectMousePressEvent()
		f.IconMaximize.Widget.DisconnectMousePressEvent()
		f.IconRestore.Widget.DisconnectMousePressEvent()
		f.IconClose.Widget.DisconnectMousePressEvent()
		f.IconMinimize.Widget.DisconnectMouseReleaseEvent()
		f.IconMaximize.Widget.DisconnectMouseReleaseEvent()
		f.IconRestore.Widget.DisconnectMouseReleaseEvent()
		f.IconClose.Widget.DisconnectMouseReleaseEvent()
		t.DisconnectMousePressEvent()
		t.DisconnectMouseReleaseEvent()
		t.DisconnectMouseMoveEvent()
		t.DisconnectMouseDoubleClickEvent()

		f.IconMinimize.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconMinimize.buttonColorChangeForLinux(hoverColor, "minimize")
		})
		f.IconMaximize.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconMaximize.buttonColorChangeForLinux(hoverColor, "maximize")
		})
		f.IconRestore.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconRestore.buttonColorChangeForLinux(hoverColor, "restore")
		})
		f.IconClose.Widget.ConnectEnterEvent(func(event *core.QEvent) {
			f.IconClose.buttonColorChangeForLinux(hoverCloseColor, "close")
		})

		f.IconMinimize.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconMinimize.buttonColorChangeForLinux(color, "minimize")
		})
		f.IconMaximize.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconMaximize.buttonColorChangeForLinux(color, "maximize")
		})
		f.IconRestore.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconRestore.buttonColorChangeForLinux(color, "restore")
		})
		f.IconClose.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
			f.IconClose.buttonColorChangeForLinux(closeColor, "close")
		})
	}

	// Button Actions
	f.IconMinimize.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.IconMaximize.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.IconRestore.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.IconClose.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.IconMinimize.Widget.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		isContain := f.IconMinimize.Widget.Rect().Contains(e.Pos(), false)
		if !isContain {
			return
		}
		f.WindowMinimize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconMaximize.Widget.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		isContain := f.IconMinimize.Widget.Rect().Contains(e.Pos(), false)
		if !isContain {
			return
		}
		f.WindowMaximize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconRestore.Widget.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		isContain := f.IconMinimize.Widget.Rect().Contains(e.Pos(), false)
		if !isContain {
			return
		}
		f.WindowRestore()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconClose.Widget.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		isContain := f.IconMinimize.Widget.Rect().Contains(e.Pos(), false)
		if !isContain {
			return
		}
		f.Close()
	})

	// TitleBar Actions
	t.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Raise()
		f.IsTitleBarPressed = true
		f.TitleBarMousePos = e.GlobalPos()
		f.Position = f.Pos()
	})

	t.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	t.ConnectMouseMoveEvent(func(e *gui.QMouseEvent) {
		if !f.IsTitleBarPressed {
			return
		}
		x := f.Position.X() + e.GlobalPos().X() - f.TitleBarMousePos.X()
		y := f.Position.Y() + e.GlobalPos().Y() - f.TitleBarMousePos.Y()
		newPos := core.NewQPoint2(x, y)
		f.Move(newPos)
	})

	t.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
		if f.IconMaximize.Widget.IsVisible() {
			f.WindowMaximize()
		} else {
			f.WindowRestore()
		}
	})
}


func (f *QToolButtonForNotDarwin) buttonColorChangeForLinux(color *RGB, buttonType string) {
	var svg string

	switch buttonType {
	case "minimize":
		svg = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M17,13H7V11H17M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
		</svg>
		`, color.Hex())
	case "maximize":
		svg = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
		<g transform="scale(0.6) translate(8,8)">
			<path fill="%s" d="M19,3H5C3.89,3 3,3.89 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19V5C21,3.89 20.1,3 19,3M19,5V19H5V5H19Z" />
		</g>
		</svg>
		`, color.Hex(), f.f.WindowColor.Hex())
	case "restore":
		svg = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<path fill="%s" d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
		<g transform="scale(0.6) translate(8,8)">
			<path fill="%s" d="M19,3H5C3.89,3 3,3.89 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19V5C21,3.89 20.1,3 19,3M19,5V19H5V5H19Z" />
		</g>
		</svg>
		`, color.Hex(), f.f.WindowColor.Hex())
	case "close":
		svg = fmt.Sprintf(`
		<svg style="width:24px;height:24px" viewBox="0 0 24 24">
		<g transform="translate(0,1)">
		<path fill="%s" d="M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"/><path d="M0 0h24v24H0z" fill="none"/></g></svg>
		`, color.Hex())
	}

	f.IconBtn.Load2(core.NewQByteArray2(svg, len(svg)))
}

func (f *QFramelessWindow) WindowMinimize() {
	f.SetWindowState(core.Qt__WindowMinimized)
}

func (f *QFramelessWindow) WindowMaximize() {
	f.IconMaximize.Widget.SetVisible(false)
	f.IconRestore.Widget.SetVisible(true)
	f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.SetWindowState(core.Qt__WindowFullScreen)
	f.IconRestore.SetStyle(nil)
}

func (f *QFramelessWindow) WindowRestore() {
	f.IconMaximize.Widget.SetVisible(true)
	f.IconRestore.Widget.SetVisible(false)
	f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.SetWindowState(core.Qt__WindowNoState)
}
