package qframelesswindow

import (
	"unsafe"

	"github.com/therecipe/qt/core"
	win "github.com/akiyosi/w32"
)

func (f *QFramelessWindow) SetNativeEvent() {
	f.winid = f.Window.WinId()
	f.Window.ConnectNativeEvent(func(eventType *core.QByteArray, message unsafe.Pointer, result *int) bool {
		msg := (*win.MSG)(message)
		hwnd := msg.Hwnd
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
	
		case win.WM_NCACTIVATE:
			f.putShadow(hwnd)

		}
		return false
	})
}

func (f *QFramelessWindow) putShadow(hwnd win.HWND) {
	if f.borderless {
		return
	}
	// style
	style := win.GetWindowLong(hwnd, win.GWL_STYLE)
	style = style | win.WS_THICKFRAME
	styleptr := uintptr(unsafe.Pointer(&style))
	ret1 := win.SetWindowLongPtr(hwnd, win.GWL_STYLE, styleptr)
	if ret1 == 0 {
		return
	}

	// shadow
	shadow := &win.MARGINS{1, 1, 1, 1}
	ret2 := win.DwmExtendFrameIntoClientArea(hwnd, shadow)
	if ret2 != 0 {
		return
	}

	var uflag uint
	uflag = win.SWP_NOZORDER | win.SWP_NOOWNERZORDER | win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_FRAMECHANGED
	var nullptr win.HWND
	ret3 := win.SetWindowPos(hwnd, nullptr, 0, 0, 0, 0, uflag)
	if !ret3 {
		return
	}
	win.ShowWindow(hwnd, win.SW_SHOW)
	f.winid = uintptr(hwnd)
	f.borderless = true
}
