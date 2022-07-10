package qframelesswindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setNSWindowStyle(long *wid, bool isVisibleTitlebar, bool isTransparent, bool isFullscreen) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    // Style
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    window.styleMask |= NSWindowStyleMaskTitled;
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskClosable;

    [[window standardWindowButton:NSWindowCloseButton] setEnabled:YES];

    if (isVisibleTitlebar) {
        [[window standardWindowButton:NSWindowCloseButton] setHidden:NO];
        [[window standardWindowButton:NSWindowMiniaturizeButton] setHidden:NO];
        [[window standardWindowButton:NSWindowZoomButton] setHidden:NO];
    } else {
        [[window standardWindowButton:NSWindowCloseButton] setHidden:YES];
        [[window standardWindowButton:NSWindowMiniaturizeButton] setHidden:YES];
        [[window standardWindowButton:NSWindowZoomButton] setHidden:YES];
    }

    // Don't show title bar
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;

    // Appearance
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];

    if (!isTransparent) {
        window.hasShadow = YES;
    } else {
        window.hasShadow = NO;
    }

    // Move buttons position when fullscreen
    if (!isFullscreen) {
        CGFloat x = 12;
        CGFloat y = -2;
        [[window standardWindowButton:NSWindowCloseButton] setFrameOrigin:NSMakePoint(x, y)];
        x += 20;
        [[window standardWindowButton:NSWindowMiniaturizeButton] setFrameOrigin:NSMakePoint(x, y)];
        x += 20;
        [[window standardWindowButton:NSWindowZoomButton] setFrameOrigin:NSMakePoint(x, y)];
    }
}


*/
import "C"

import (
	"unsafe"

	"github.com/therecipe/qt/core"
)

func (f *QFramelessWindow) SetStyleMask() {
	f.SetNSWindowStyleMask(
		!f.IsTitlebarHidden,
		f.WindowColorAlpha != 1.0,
		f.WindowState() == core.Qt__WindowFullScreen,
	)
}

func (f *QFramelessWindow) SetNSWindowStyleMask(isVisibleTitlebarButtons, hasAlpha, isWindowFullscreen bool) {
	wid := f.WinId()
	C.setNSWindowStyle(
		(*C.long)(unsafe.Pointer(wid)),
		C.bool(isVisibleTitlebarButtons),
		C.bool(hasAlpha),
		C.bool(isWindowFullscreen),
	)
}

func (f *QFramelessWindow) SetupNativeEvent() {
}

func (f *QFramelessWindow) SetupNativeEvent2() {
}
