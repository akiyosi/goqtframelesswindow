package qframelesswindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setStyleMask(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    window.styleMask |= NSResizableWindowMask;
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.movableByWindowBackground = YES;

    return;
}
*/
import "C"

import (
	"unsafe"
)

func (f *QFramelessWindow) SetStyleMask() {
	wid := f.WinId()
	C.setStyleMask((*_Ctype_long)(unsafe.Pointer(wid)))
}

func (f *QFramelessWindow) SetupNativeEvent() {
	// f.ConnectNativeEvent(func(eventType *core.QByteArray, message unsafe.Pointer, result *int) bool {
	// 	fmt.Println("msg", message)

	// 	return false
	// })
}
