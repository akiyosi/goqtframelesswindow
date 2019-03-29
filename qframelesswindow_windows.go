package qframelesswindow

import (
	"unsafe"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"

	win "github.com/akiyosi/w32"
)

func (f *QFramelessWindow) SetNativeEvent(app *widgets.QApplication) {
	filterObj := core.NewQAbstractNativeEventFilter()
	filterObj.ConnectNativeEventFilter(func(eventType *core.QByteArray, message unsafe.Pointer, result int) bool {
		msg := (*win.MSG)(message)
		hwnd := msg.Hwnd

		switch msg.Message {
		case win.WM_CREATE:
		 	style := win.GetWindowLong(hwnd, win.GWL_STYLE)
		 	style = style | win.WS_THICKFRAME | win.WS_CAPTION
		 	win.SetWindowLong(hwnd, win.GWL_STYLE, uint32(style))
			shadow := &win.MARGINS{-1, -1, -1, -1}
			win.DwmExtendFrameIntoClientArea(hwnd, shadow)
			f.borderless = true
		 	return false

		case win.WM_NCCALCSIZE:
			if msg.WParam == 1 && f.borderless {
				// this kills the window frame and title bar we added with WS_THICKFRAME and WS_CAPTION
				result = 0
				if win.IsWindowVisible(hwnd) == false {
		 			win.ShowWindow(hwnd, win.SW_SHOW)
				}
			}
			return true

		default:
		}
		return false
	})
	app.InstallNativeEventFilter(filterObj)
}
