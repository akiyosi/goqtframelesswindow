package qframelesswindow

import (
	"fmt"
	"unsafe"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func (f *QFramelessWindow) SetNativeEvent(app *widgets.QApplication) {
	filterObj := core.NewQAbstractNativeEventFilter()
	filterObj.ConnectNativeEventFilter(func(eventType *core.QByteArray, message unsafe.Pointer, result int) bool {
		fmt.Println("debug:", eventType)

		// return filterObj.NativeEventFilter(eventType, message, result)
		return false
	})
	app.InstallNativeEventFilter(filterObj)
}
