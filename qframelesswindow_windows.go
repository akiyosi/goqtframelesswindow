package qframelesswindow

import (
	"unsafe"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/therecipe/qt/gui"

	win "github.com/akiyosi/w32"
)

func (f *QFramelessWindow) SetNativeEvent(app *widgets.QApplication) {
	filterObj := core.NewQAbstractNativeEventFilter()
	filterObj.ConnectNativeEventFilter(func(eventType *core.QByteArray, message unsafe.Pointer, result int) bool {
		msg := (*win.MSG)(message)
		lparam := msg.LParam
		hwnd := msg.Hwnd

		switch msg.Message {
		case win.WM_CREATE:
		 	style := win.GetWindowLong(hwnd, win.GWL_STYLE)
		 	style = style | win.WS_THICKFRAME | win.WS_CAPTION
		 	win.SetWindowLong(hwnd, win.GWL_STYLE, uint32(style))
			f.borderless = true
		 	return false

		case win.WM_NCCALCSIZE:
			if msg.WParam == 1 && f.borderless {
				// this kills the window frame and title bar we added with WS_THICKFRAME and WS_CAPTION
				result = 0
				if win.IsWindowVisible(hwnd) == false {
					shadow := &win.MARGINS{1, 1, 1, 1}
					win.DwmExtendFrameIntoClientArea(hwnd, shadow)
		 			win.ShowWindow(hwnd, win.SW_SHOW)
				}
			}
			return true

		case win.WM_GETMINMAXINFO:
			mm := (*win.MINMAXINFO)((unsafe.Pointer)(lparam))
			mm.PtMinTrackSize.X = int32(f.minimumWidth)
			mm.PtMinTrackSize.Y = int32(f.minimumHeight)

			return true

		case win.WM_NCHITTEST:
			result = 0

			rect := win.GetWindowRect(hwnd)

			// Get the cursor position from Qt instead of winapi.
			// x, y, _ := win.GetCursorPos()
			x := f.MousePos[0]
			y := f.MousePos[1]

			rectX := int(rect.Left)
			rectY := int(rect.Top)
			rectWidth := int(rect.Right - rect.Left)
			rectHeight := int(rect.Bottom - rect.Top)

			edge := f.detectEdgeOnCursor(x, y, rectX, rectY, rectWidth, rectHeight)
			if !f.isLeftButtonPressed {
				f.hoverEdge = edge
			}

			if !f.isLeftButtonPressed {
				f.isCursorChanged = true
				cursor := gui.NewQCursor()
				if f.Window.IsFullScreen() || f.Window.IsMaximized() {
					if f.isCursorChanged {
						cursor.SetShape(core.Qt__ArrowCursor)
						f.Window.SetCursor(cursor)
					}
				}
				switch f.hoverEdge {
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

		default:
		}
		return false
	})
	app.InstallNativeEventFilter(filterObj)
}
