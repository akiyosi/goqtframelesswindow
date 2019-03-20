// +build !darwin

package qframelesswindow

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

func (f *QFramelessWindow) SetTitleBarActions() {
	t := f.TitleBar

	f.IconMinimize.Widget.ConnectEnterEvent(func(event *core.QEvent) {
		f.IconMinimize.SetStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})
	f.IconMaximize.Widget.ConnectEnterEvent(func(event *core.QEvent) {
		f.IconMaximize.SetStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})
	f.IconRestore.Widget.ConnectEnterEvent(func(event *core.QEvent) {
		f.IconRestore.SetStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})
	f.IconClose.Widget.ConnectEnterEvent(func(event *core.QEvent) {
		f.IconClose.SetStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})

	f.IconMinimize.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
		f.IconMinimize.SetStyle(nil)
	})
	f.IconMaximize.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
		f.IconMaximize.SetStyle(nil)
	})
	f.IconRestore.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
		f.IconRestore.SetStyle(nil)
	})
	f.IconClose.Widget.ConnectLeaveEvent(func(event *core.QEvent) {
		f.IconClose.SetStyle(nil)
	})

	// Button Actions
	f.IconMinimize.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Window.SetWindowState(core.Qt__WindowMinimized)
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconMaximize.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowMaximize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconRestore.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowRestore()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconClose.Widget.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
	})

	// TitleBar Actions
	t.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Raise()
		f.IsMousePressed = true
		f.MousePos = e.GlobalPos()
		f.Pos = f.Window.Pos()
	})

	t.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.IsMousePressed = false
	})

	t.ConnectMouseMoveEvent(func(e *gui.QMouseEvent) {
		if !f.IsMousePressed {
			return
		}
		x := f.Pos.X() + e.GlobalPos().X() - f.MousePos.X()
		y := f.Pos.Y() + e.GlobalPos().Y() - f.MousePos.Y()
		newPos := core.NewQPoint2(x, y)
		f.Window.Move(newPos)
	})

	t.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
		if f.IconMaximize.Widget.IsVisible() {
			f.windowMaximize()
		} else {
			f.windowRestore()
		}
	})
}

func (f *QFramelessWindow) windowMaximize() {
	f.IconMaximize.Widget.SetVisible(false)
	f.IconRestore.Widget.SetVisible(true)
	f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.Window.SetWindowState(core.Qt__WindowMaximized)
}

func (f *QFramelessWindow) windowRestore() {
	f.IconMaximize.Widget.SetVisible(true)
	f.IconRestore.Widget.SetVisible(false)
	f.Layout.SetContentsMargins(f.shadowMargin, f.shadowMargin, f.shadowMargin, f.shadowMargin)
	f.Window.SetWindowState(core.Qt__WindowNoState)
}
