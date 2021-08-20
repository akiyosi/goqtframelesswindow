package qframelesswindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setStyleMask(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    // Style
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;

    // Appearance
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = NO;

    return;
}

void setStyleMaskAndShadow(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    // Style
    window.styleMask |= NSWindowStyleMaskTitled;
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskClosable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;

    // Appearance
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = YES;

    // Don't show title bar
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;

    // Hidden native buttons
    [[window standardWindowButton:NSWindowCloseButton] setHidden:YES];
    [[window standardWindowButton:NSWindowMiniaturizeButton] setHidden:YES];
    [[window standardWindowButton:NSWindowZoomButton] setHidden:YES];

    return;
}

void setStyleMaskWithNativeButtons(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    // Style
    window.styleMask |= NSWindowStyleMaskTitled;
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskClosable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;

    // Appearance
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = NO;

    // Don't show title bar
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;

    // Enable close button
    [[window standardWindowButton:NSWindowCloseButton] setEnabled:YES];

    return;
}

void setStyleMaskAndShadowWithNativeButtons(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    // Style
    window.styleMask |= NSWindowStyleMaskTitled;
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskClosable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;

    // Appearance
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = YES;

    // Don't show title bar
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;

    // Enable close button
    [[window standardWindowButton:NSWindowCloseButton] setEnabled:YES];

    return;
}

void setStyleMaskWithNativeButtonsWithMove(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    // Style
    window.styleMask |= NSWindowStyleMaskTitled;
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskClosable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;

    // Appearance
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = NO;

    // Don't show title bar
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;

    // Enable close button
    [[window standardWindowButton:NSWindowCloseButton] setEnabled:YES];

    // Move position
    CGFloat x = 12;
    CGFloat y = -2;
    [[window standardWindowButton:NSWindowCloseButton] setFrameOrigin:NSMakePoint(x, y)];
    x += 20;
    [[window standardWindowButton:NSWindowMiniaturizeButton] setFrameOrigin:NSMakePoint(x, y)];
    x += 20;
    [[window standardWindowButton:NSWindowZoomButton] setFrameOrigin:NSMakePoint(x, y)];

    return;
}

void setStyleMaskAndShadowWithNativeButtonsWithMove(long *wid) {
    NSView* view = (NSView*)wid;
    NSWindow *window = view.window;

    // Style
    window.styleMask |= NSWindowStyleMaskTitled;
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskClosable;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;

    // Appearance
    window.opaque = NO;
    window.backgroundColor = [NSColor clearColor];
    window.hasShadow = YES;

    // Don't show title bar
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;

    // Enable close button
    [[window standardWindowButton:NSWindowCloseButton] setEnabled:YES];

    // Move position
    CGFloat x = 12;
    CGFloat y = -2;
    [[window standardWindowButton:NSWindowCloseButton] setFrameOrigin:NSMakePoint(x, y)];
    x += 20;
    [[window standardWindowButton:NSWindowMiniaturizeButton] setFrameOrigin:NSMakePoint(x, y)];
    x += 20;
    [[window standardWindowButton:NSWindowZoomButton] setFrameOrigin:NSMakePoint(x, y)];

    return;
}

*/
import "C"

import (
	"unsafe"

	"github.com/therecipe/qt/core"
)

func (f *QFramelessWindow) SetStyleMask() {
	wid := f.WinId()
	if f.WindowColorAlpha == 1.0 {
		if f.WindowState() == core.Qt__WindowFullScreen {
			C.setStyleMaskAndShadowWithNativeButtons((*C.long)(unsafe.Pointer(wid)))
		} else {
			C.setStyleMaskAndShadowWithNativeButtonsWithMove((*C.long)(unsafe.Pointer(wid)))
		}
	} else {
		if f.WindowState() == core.Qt__WindowFullScreen {
			C.setStyleMaskWithNativeButtons((*C.long)(unsafe.Pointer(wid)))
		} else {
			C.setStyleMaskWithNativeButtonsWithMove((*C.long)(unsafe.Pointer(wid)))
		}
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
