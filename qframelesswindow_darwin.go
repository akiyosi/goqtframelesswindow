package qframelesswindow

import (
	"fmt"
)

func (f *QFramelessWindow) SetupNativeEvent() {
	fmt.Println("Need to imprement SetNativeEvent() for darwin")
	// filterObj := core.NewQAbstractNativeEventFilter()
	// filterObj.ConnectNativeEventFilter(func(eventType *core.QByteArray, message unsafe.Pointer, result int) bool {

	// 	// return filterObj.NativeEventFilter(eventType, message, result)
	// 	return false
	// })
	// app.InstallNativeEventFilter(filterObj)
}
