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
		msg := (*win.MSG)(message)
		switch msg.Message {
		case win.WM_NCCALCSIZE:
			fmt.Println("debug")
			winid := (*win.HWND)(unsafe.Pointer(f.Window.WinId()))
			style := win.GetWindowLong(*winid, win.GWL_STYLE)
			style = style | win.WS_MAXIMIZEBOX | win.WS_THICKFRAME | win.WS_CAPTION
			win.SetWindowLong(*winid, win.GWL_STYLE, style)
			// case win.WM_NCHITTEST:
			// 	var winrect *win.RECT
			//  	winid := (*win.HWND)(unsafe.Pointer(f.Window.WinId()))
			// 	win.GetWindowRect(*winid, winrect)
			// 	fmt.Println(winrect.Left, winrect.Bottom)
		}
		return filterObj.NativeEventFilter(eventType, message, result)
	})
	app.InstallNativeEventFilter(filterObj)
}
