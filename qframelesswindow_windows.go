package qframelesswindow

import (
	"unsafe"

	"github.com/therecipe/qt/core"
	win "github.com/akiyosi/w32"
)

func (f *QFramelessWindow) SetNativeEvent() {
	f.Window.WinId()
	f.Window.ConnectNativeEvent(func(eventType *core.QByteArray, message unsafe.Pointer, result *int) bool {
		msg := (*win.MSG)(message)
		hwnd := msg.Hwnd

		switch msg.Message {
		case win.WM_NCCALCSIZE:
			if msg.WParam == 1 {
				*result = 0
				return true
			}
			*result = (int)(win.DefWindowProc(msg.Hwnd, win.WM_NCCALCSIZE, msg.WParam, msg.LParam))
			return true

		case win.WM_GETMINMAXINFO:
			mm := (*win.MINMAXINFO)((unsafe.Pointer)(msg.LParam))
			mm.PtMinTrackSize.X = int32(f.minimumWidth)
			mm.PtMinTrackSize.Y = int32(f.minimumHeight)
                
			return true
	
		case win.WM_ACTIVATE:
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
	win.SetWindowLong(hwnd, win.GWL_STYLE, uint32(style))

	// shadow
	shadow := &win.MARGINS{1, 1, 1, 1}
	win.DwmExtendFrameIntoClientArea(hwnd, shadow)

	var uflag uint
	uflag = win.SWP_NOZORDER | win.SWP_NOOWNERZORDER | win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_FRAMECHANGED
	var nullptr win.HWND
	win.SetWindowPos(hwnd, nullptr, 0, 0, 0, 0, uflag)

	win.UpdateWindow(hwnd)

	f.borderless = true
}
