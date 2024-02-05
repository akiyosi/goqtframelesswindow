package qframelesswindow

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>


void applyBlurEffect(long *wid, long *contentid, bool isLight) {
    NSView* view = (NSView*)wid;
    NSWindow *window = [view window];


    // NSVisualEffectViewをウィンドウのcontentViewに追加
    NSVisualEffectView *visualEffectView = [[NSVisualEffectView alloc] initWithFrame:window.contentView.bounds];
    visualEffectView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    visualEffectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
    visualEffectView.state = NSVisualEffectStateActive;

	if (isLight) {
		visualEffectView.material = NSVisualEffectMaterialLight;
	} else {
		visualEffectView.material = NSVisualEffectMaterialUltraDark;
	}

    [window.contentView addSubview:visualEffectView positioned:NSWindowBelow relativeTo:nil];

    NSView* qtView = (NSView*)contentid;

    // QtのビューをvisualEffectViewの上に配置
    [window.contentView addSubview:qtView positioned:NSWindowAbove relativeTo:visualEffectView];
}

void setNSWindowStyle(long *wid, bool isVisibleTitlebar, float red, float green, float blue, float alpha, bool isFullscreen) {
    NSView* view = (NSView*)wid;
    NSWindow *window = [view window];

    // Style
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    if (isVisibleTitlebar) {
        window.styleMask |= NSWindowStyleMaskTitled;
    }
    window.styleMask |= NSWindowStyleMaskResizable;
    window.styleMask |= NSWindowStyleMaskMiniaturizable;
    window.styleMask |= NSWindowStyleMaskClosable;

    [[window standardWindowButton:NSWindowCloseButton] setEnabled:YES];
    [[window standardWindowButton:NSWindowCloseButton] setHidden:!isVisibleTitlebar];
    [[window standardWindowButton:NSWindowMiniaturizeButton] setHidden:!isVisibleTitlebar];
    [[window standardWindowButton:NSWindowZoomButton] setHidden:!isVisibleTitlebar];

    // Don't show title bar
    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = isVisibleTitlebar ? NSWindowTitleVisible : NSWindowTitleHidden;

    // Appearance
    window.opaque = NO;
    CGFloat cgred = red;
    CGFloat cggreen = green;
    CGFloat cgblue = blue;
    CGFloat cgalpha = alpha;
    window.backgroundColor = [NSColor colorWithCalibratedRed:cgred green:cggreen blue:cgblue alpha:cgalpha];
    window.hasShadow = YES;

    // Fullscreen
    if (isFullscreen) {
        window.styleMask |= NSWindowStyleMaskFullScreen;
    }
}

*/
import "C"

import (
	"unsafe"

	"github.com/akiyosi/qt/core"
)

func (f *QFramelessWindow) SetStyleMask() {
	f.SetNSWindowStyleMask(
		!f.IsTitlebarHidden,
		f.WindowColor.R, f.WindowColor.G, f.WindowColor.B,
		float32(f.WindowColorAlpha),
		f.WindowState() == core.Qt__WindowFullScreen,
	)
}

func (f *QFramelessWindow) SetBlurEffectForMacOS(isLight bool) {
	wid := f.WinId()
	contentid := f.Content.WinId()
	C.applyBlurEffect(
		(*C.long)(unsafe.Pointer(wid)),
		(*C.long)(unsafe.Pointer(contentid)),
		C.bool(isLight),
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

func (f *QFramelessWindow) SetBlurEffectForWin(hwnd uintptr) {
}

func (f *QFramelessWindow) SetupNativeEvent() {
}

func (f *QFramelessWindow) SetupNativeEvent2() {
}
