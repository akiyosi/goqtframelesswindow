package qframelesswindow

import (
	"unsafe"

	"github.com/akiyosi/qt/core"
	win "github.com/akiyosi/w32"
)

func (f *QFramelessWindow) SetNSWindowStyleMask(isVisibleTitlebarButtons bool, R, G, B uint16, alpha float32, isWindowFullscreen bool) {
}

func (f *QFramelessWindow) SetStyleMask() {
}

func (f *QFramelessWindow) SetBlurEffectForMacOS(isLight bool) {
}

// SetBlurEffect function applies blur effect to the given window
func (f *QFramelessWindow) SetBlurEffectForWin(hwnd uintptr) {
	if !f.ApplyBlurEffect {
		return
	}

	var accentFlags uint32
	accentFlags = 0x20 | 0x40 | 0x80 | 0x100 // enable shadow

	R := uint32(f.WindowColor.R)
	G := uint32(f.WindowColor.G)
	B := uint32(f.WindowColor.B)
	gradientColor := uint32((0x99 << 24) | (R << 16) | (G << 8) | B)

	var accent = win.ACCENTPOLICY{
		AccentState:   win.ACCENT_ENABLE_ACRYLICBLURBEHIND,
		AccentFlags:   accentFlags,
		GradientColor: gradientColor,
		AnimationId:   uint32(0),
	}

	var data = &win.WINDOWCOMPOSITIONATTRIBDATA{
		Attribute:  19, // WCA_ACCENT_POLICY
		Data:       unsafe.Pointer(&accent),
		SizeOfData: uint32(unsafe.Sizeof(accent)),
	}

	win.SetWindowCompositionAttribute(win.HWND(hwnd), data)
}

// func (f *QFramelessWindow) SetupNativeEvent2() {
// 	f.WinId()
// 	f.ConnectNativeEvent(func(eventType *core.QByteArray, message unsafe.Pointer, result *int) bool {
// 		msg := (*win.MSG)(message)
// 		hwnd := msg.Hwnd
//
// 		switch msg.Message {
// 		case win.WM_NCCALCSIZE:
// 			if msg.WParam == 1 {
// 				*result = 0
// 				return true
// 			}
// 			*result = (int)(win.DefWindowProc(msg.Hwnd, win.WM_NCCALCSIZE, msg.WParam, msg.LParam))
// 			return true
//
// 		case win.WM_GETMINMAXINFO:
// 			mm := (*win.MINMAXINFO)((unsafe.Pointer)(msg.LParam))
// 			mm.PtMinTrackSize.X = int32(f.minimumWidth)
// 			mm.PtMinTrackSize.Y = int32(f.minimumHeight)
//
// 			return true
//
// 		case win.WM_ACTIVATEAPP:
// 			// case win.WM_NCACTIVATE:
// 			f.putTransparent(hwnd)
// 			// f.putShadow(hwnd)
// 		}
//
// 		return false
// 	})
// }

func (f *QFramelessWindow) SetupNativeEvent() {
	f.WinId()
	f.ConnectNativeEvent(func(eventType *core.QByteArray, message unsafe.Pointer, result *int) bool {
		msg := (*win.MSG)(message)
		hwnd := msg.Hwnd
		f.hwnd = uintptr(hwnd)

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

		// case win.WM_STYLECHANGING:
		// case win.WM_STYLECHANGED:
		// case win.WM_SHOWWINDOW:
		case win.WM_ACTIVATEAPP:
			// case win.WM_NCACTIVATE:
			f.putShadow(hwnd)

		}

		// return false
		return f.NativeEventDefault(eventType, message, result)
	})
}

func (f *QFramelessWindow) setLayerd(hwnd win.HWND) {
	if f.borderless {
		return
	}
	// style
	style := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)
	style = win.WS_EX_LAYERED | win.WS_EX_TRANSPARENT
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, uint32(style))
	win.SetLayeredWindowAttributes(hwnd, 0, 255, 0x00000001)
	win.UpdateWindow(hwnd)

	f.borderless = true
}

func (f *QFramelessWindow) putTransparent(hwnd win.HWND) {
	if f.borderless {
		return
	}
	// style
	style := win.GetWindowLong(hwnd, win.GWL_STYLE)
	// style = style | win.WS_THICKFRAME | win.WS_MAXIMIZEBOX | win.WS_CAPTION
	// style = style & ^(win.WS_CAPTION | win.WS_HSCROLL | win.WS_VSCROLL | win.WS_SYSMENU | win.WS_MAXIMIZE)
	style = style | win.WS_THICKFRAME | win.WS_MAXIMIZEBOX | win.WS_CAPTION
	win.SetWindowLong(hwnd, win.GWL_STYLE, uint32(style))

	// // exstyle := win.WS_EX_LAYERED | win.WS_EX_TRANSPARENT
	exstyle := win.WS_EX_LAYERED
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, uint32(exstyle))
	// alpha := uint8(math.Sqrt(f.WindowColorAlpha) * 255)
	// black := &RGB{0, 0, 0}
	// win.SetLayeredWindowAttributes(hwnd, black.toColorref(), 0, win.LWA_COLORKEY)
	// win.SetLayeredWindowAttributes(hwnd, 0, 10, win.LWA_ALPHA)

	// shadow
	shadow := &win.MARGINS{-1, -1, -1, -1}
	win.DwmExtendFrameIntoClientArea(hwnd, shadow)

	var uflag uint
	uflag = win.SWP_NOZORDER | win.SWP_NOOWNERZORDER | win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_FRAMECHANGED
	// uflag = win.SWP_SHOWWINDOW | win.SWP_FRAMECHANGED
	var nullptr win.HWND
	win.SetWindowPos(hwnd, nullptr, 0, 0, 0, 0, uflag)

	win.UpdateWindow(hwnd)

	f.borderless = true
}

func (f *QFramelessWindow) putShadow(hwnd win.HWND) {
	if f.borderless {
		return
	}
	// style
	style := win.GetWindowLong(hwnd, win.GWL_STYLE)
	style = style | win.WS_THICKFRAME | win.WS_MAXIMIZEBOX | win.WS_CAPTION
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
