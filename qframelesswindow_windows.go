package qframelesswindow

import (
	"fmt"
	"unsafe"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"

	"github.com/lxn/win"
)

func (f *QFramelessWindow) SetNativeEvent(app *widgets.QApplication) {
	filterObj := core.NewQAbstractNativeEventFilter()
	filterObj.ConnectNativeEventFilter(func(eventType *core.QByteArray, message unsafe.Pointer, result int) bool {
		fmt.Println(message)
		msg := (*win.MSG)(message)
		switch msg.Message {
		case win.WM_KEYDOWN:
			return false
		case win.WM_NCCALCSIZE:
			fmt.Println("debug")
			winid := (*win.HWND)(unsafe.Pointer(f.Window.WinId()))
			style := win.GetWindowLong(*winid, win.GWL_STYLE)
			style = style | win.WS_MAXIMIZEBOX | win.WS_THICKFRAME | win.WS_CAPTION
			win.SetWindowLong(*winid, win.GWL_STYLE, style)
			return true
		case win.WM_NCHITTEST:
			fmt.Println("debug:: WM_NCHITTEST")
			var winrect *win.RECT
			winid := (*win.HWND)(unsafe.Pointer(f.Window.WinId()))
			win.GetWindowRect(*winid, winrect)
			fmt.Println(winrect.Left, winrect.Bottom)
			return true
		default:
			fmt.Println("debug --", msg.Message)
			return true
		}

		// return filterObj.NativeEventFilter(eventType, message, result)
	})
	app.InstallNativeEventFilter(filterObj)
}
