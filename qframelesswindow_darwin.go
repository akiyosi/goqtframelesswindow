package qframelesswindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setNSWindowStyle(long *wid, bool isVisibleTitlebar, float red, float green, float blue, float alpha, bool isFullscreen) {
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
    CGFloat cgred = red;
    CGFloat cggreen = green;
    CGFloat cgblue = blue;
    CGFloat cgalpha = alpha;
    window.alphaValue = 1.0;
    // window.backgroundColor = [NSColor clearColor];
    window.backgroundColor = [NSColor colorWithCalibratedRed:cgred green:cggreen blue:cgblue alpha:cgalpha];
    window.hasShadow = YES;


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
		f.WindowColor.R, f.WindowColor.G, f.WindowColor.B,
		float32(f.WindowColorAlpha),
		f.WindowState() == core.Qt__WindowFullScreen,
	)
}

func (f *QFramelessWindow) SetNSWindowStyleMask(isVisibleTitlebarButtons bool, R, G, B uint16, alpha float32, isWindowFullscreen bool) {
	wid := f.WinId()
	fR := float32(R) / float32(255)
	fG := float32(G) / float32(255)
	fB := float32(B) / float32(255)
	C.setNSWindowStyle(
		(*C.long)(unsafe.Pointer(wid)),
		C.bool(isVisibleTitlebarButtons),
		C.float(fR), C.float(fG), C.float(fB),
		C.float(alpha),
		C.bool(isWindowFullscreen),
	)
}

func (f *QFramelessWindow) SetupNativeEvent() {
}

func (f *QFramelessWindow) SetupNativeEvent2() {
}
