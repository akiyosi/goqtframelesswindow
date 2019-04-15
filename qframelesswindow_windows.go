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
		lparam := msg.LParam
		hwnd := msg.Hwnd
		var uflag uint
		uflag = win.SWP_NOZORDER | win.SWP_NOOWNERZORDER | win.SWP_NOMOVE | win.SWP_NOSIZE | win.SWP_FRAMECHANGED
		var nullptr win.HWND
		shadow := &win.MARGINS{0, 0, 0, 1}

		switch msg.Message {
		case win.WM_CREATE:
			// set style
		 	// style := win.GetWindowLong(hwnd, win.GWL_STYLE)
		 	// style = style | win.WS_THICKFRAME | win.WS_CAPTION
			style := win.WS_POPUP | win.WS_THICKFRAME | win.WS_MINIMIZEBOX | win.WS_MAXIMIZEBOX | win.WS_CAPTION
		 	win.SetWindowLong(hwnd, win.GWL_STYLE, uint32(style))

			// shadow
			win.DwmExtendFrameIntoClientArea(hwnd, shadow)
			win.SetWindowPos(hwnd, nullptr, 0, 0, 0, 0, uflag)

		 	return true

		case win.WM_NCCALCSIZE:
			if msg.WParam == 1 {
				// this kills the window frame and title bar we added with WS_THICKFRAME and WS_CAPTION
				result = 0
				return true
			}
			return false

		case win.WM_GETMINMAXINFO:
			mm := (*win.MINMAXINFO)((unsafe.Pointer)(lparam))
			mm.PtMinTrackSize.X = int32(f.minimumWidth)
			mm.PtMinTrackSize.Y = int32(f.minimumHeight)

			return true

	// 	case win.WM_ACTIVATEAPP:
	// 		fmt.Println("win.WM_ACTIVATEAPP:")

	// 	case win.WM_CANCELMODE:
	// 		fmt.Println("win.WM_CANCELMODE:")

	// 	case win.WM_CHILDACTIVATE:
	// 		fmt.Println("win.WM_CHILDACTIVATE:")

	// 	case win.WM_CLOSE:
	// 		fmt.Println("win.WM_CLOSE:")

	// 	case win.WM_COMPACTING:
	// 		fmt.Println("win.WM_COMPACTING:")

	// 	case win.WM_DESTROY:
	// 		fmt.Println("win.WM_DESTROY:")

	// 	case win.WM_ENABLE:
	// 		fmt.Println("win.WM_ENABLE:")

	// 	case win.WM_ENTERSIZEMOVE:
	// 		fmt.Println("win.WM_ENTERSIZEMOVE:")

	// 	case win.WM_EXITSIZEMOVE:
	// 		fmt.Println("win.WM_EXITSIZEMOVE:")

	// 	case win.WM_GETICON:
	// 		fmt.Println("win.WM_GETICON:")

	// 	case win.WM_INPUTLANGCHANGE:
	// 		fmt.Println("win.WM_INPUTLANGCHANGE:")

	// 	case win.WM_INPUTLANGCHANGEREQUEST:
	// 		fmt.Println("win.WM_INPUTLANGCHANGEREQUEST:")

	// 	case win.WM_MOVE:
	// 		fmt.Println("win.WM_MOVE:")

	// 	case win.WM_MOVING:
	// 		fmt.Println("win.WM_MOVING:")

	// 	case win.WM_NCACTIVATE:
	// 		fmt.Println("win.WM_NCACTIVATE:")

	// 	case win.WM_NCCREATE:
	// 		fmt.Println("win.WM_NCCREATE:")

	// 	case win.WM_NCDESTROY:
	// 		fmt.Println("win.WM_NCDESTROY:")

	// 	case win.WM_NULL:
	// 		fmt.Println("win.WM_NULL:")

	// 	case win.WM_QUERYDRAGICON:
	// 		fmt.Println("win.WM_QUERYDRAGICON:")

	// 	case win.WM_QUERYOPEN:
	// 		fmt.Println("win.WM_QUERYOPEN:")

	// 	case win.WM_QUIT:
	// 		fmt.Println("win.WM_QUIT:")

	// 	case win.WM_SHOWWINDOW:
	// 		fmt.Println("win.WM_SHOWWINDOW:")

	// 	case win.WM_SIZE:
	// 		fmt.Println("win.WM_SIZE:")

	// 	case win.WM_SIZING:
	// 		fmt.Println("win.WM_SIZING:")

	// 	case win.WM_STYLECHANGED:
	// 		fmt.Println("win.WM_STYLECHANGED:")

	// 	case win.WM_STYLECHANGING:
	// 		fmt.Println("win.WM_STYLECHANGING:")

	// 	case win.WM_THEMECHANGED:
	// 		fmt.Println("win.WM_THEMECHANGED:")

	// 	case win.WM_USERCHANGED:
	// 		fmt.Println("win.WM_USERCHANGED:")

	// 	case win.WM_WINDOWPOSCHANGED:
	// 		fmt.Println("win.WM_WINDOWPOSCHANGED:")

	// 	case win.WM_WINDOWPOSCHANGING:
	// 		fmt.Println("win.WM_WINDOWPOSCHANGING:")

		default:
		}
		return false
	})
	app.InstallNativeEventFilter(filterObj)
}
