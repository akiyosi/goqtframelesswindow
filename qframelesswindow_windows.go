package qframelesswindow

import (
	"unsafe"

	"github.com/therecipe/qt/core"
	win "github.com/akiyosi/w32"
)


func (f *QFramelessWindow) SetNativeEvent() {
	f.Window.ConnectNativeEvent(func(eventType *core.QByteArray, message unsafe.Pointer, result *int) bool {
		msg := (*win.MSG)(message)
		hwnd := (win.HWND)(f.winid)
		lparam := msg.LParam

		switch msg.Message {
		case win.WM_NCCALCSIZE:
			if msg.WParam == 1 {
				win.SetWindowLong(hwnd, win.DWL_MSGRESULT, 0)
				return true
			}
			return false

		case win.WM_GETMINMAXINFO:
			mm := (*win.MINMAXINFO)((unsafe.Pointer)(lparam))
			mm.PtMinTrackSize.X = int32(f.minimumWidth)
			mm.PtMinTrackSize.Y = int32(f.minimumHeight)
                
			return true

		case win.WM_ACTIVATE:
			if f.borderless {
				return false
			}
			// style
			style := win.GetWindowLong(hwnd, win.GWL_STYLE)
			style = style | win.WS_THICKFRAME
			styleptr := uintptr(unsafe.Pointer(&style))
			win.SetWindowLongPtr(hwnd, win.GWL_STYLE, styleptr)

			// shadow
			shadow := &win.MARGINS{1, 1, 1, 1}
			win.DwmExtendFrameIntoClientArea(hwnd, shadow)

			win.ShowWindow(hwnd, win.SW_SHOW)
			f.borderless = true
		}
		return false
	})
}
