package qframelesswindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setStyleMask(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = NO;

    return;
}

void setStyleMaskAndShadow(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    window.styleMask |= NSWindowStyleMaskTitled;
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = YES;
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;

    [[window standardWindowButton:NSWindowCloseButton] setHidden:YES];
    [[window standardWindowButton:NSWindowMiniaturizeButton] setHidden:YES];
    [[window standardWindowButton:NSWindowZoomButton] setHidden:YES];

    return;
}
*/
import "C"

import (
	"unsafe"
)

func (f *QFramelessWindow) SetStyleMask() {
	wid := f.WinId()
	if f.WindowColorAlpha == 1.0 {
		C.setStyleMaskAndShadow((*C.long)(unsafe.Pointer(wid)))
	} else {
		C.setStyleMask((*C.long)(unsafe.Pointer(wid)))
	}
}

func (f *QFramelessWindow) SetupNativeEvent() {
	// f.ConnectNativeEvent(func(eventType *core.QByteArray, message unsafe.Pointer, result *int) bool {
	// 	fmt.Println("msg", message)

	// 	return false
	// })
}

func (f *QFramelessWindow) SetupNativeEvent2() {
}
