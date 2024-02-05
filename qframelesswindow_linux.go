package qframelesswindow

// import (
// 	"fmt"
// )

func (f *QFramelessWindow) SetNSWindowStyleMask(isVisibleTitlebarButtons bool, R, G, B uint16, alpha float32, isWindowFullscreen bool) {
}

func (f *QFramelessWindow) SetStyleMask() {
}

func (f *QFramelessWindow) SetBlurEffectForMacOS(isLight bool) {
}

func (f *QFramelessWindow) SetBlurEffectForWin(hwnd uintptr) {
}

func (f *QFramelessWindow) SetupNativeEvent() {
	// filterObj := core.NewQAbstractNativeEventFilter()
	// filterObj.ConnectNativeEventFilter(func(eventType *core.QByteArray, message unsafe.Pointer, result int) bool {

	// 	// return filterObj.NativeEventFilter(eventType, message, result)
	// 	return false
	// })
	// app.InstallNativeEventFilter(filterObj)
}

func (f *QFramelessWindow) SetupNativeEvent2() {
}
