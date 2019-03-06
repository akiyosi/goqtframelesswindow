// +build !darwin

package qframelesswindow

import (
        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
)

func (f *QFramelessWindow) SetTitleBarActions() {
	t := f.TitleBar

	f.IconMinimize.ConnectEnterEvent(func(event *core.QEvent) {
		f.SetIconMinimizeStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})
	f.IconMaximize.ConnectEnterEvent(func(event *core.QEvent) {
		f.SetIconMaximizeStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})
	f.IconRestore.ConnectEnterEvent(func(event *core.QEvent) {
		f.SetIconRestoreStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})
	f.IconClose.ConnectEnterEvent(func(event *core.QEvent) {
		f.SetIconCloseStyle(&RGB{
			R: 0,
			G: 162,
			B: 232,
		})
	})

	f.IconMinimize.ConnectLeaveEvent(func(event *core.QEvent) {
		f.SetIconMinimizeStyle(nil)
	})
	f.IconMaximize.ConnectLeaveEvent(func(event *core.QEvent) {
		f.SetIconMaximizeStyle(nil)
	})
	f.IconRestore.ConnectLeaveEvent(func(event *core.QEvent) {
		f.SetIconRestoreStyle(nil)
	})
	f.IconClose.ConnectLeaveEvent(func(event *core.QEvent) {
		f.SetIconCloseStyle(nil)
	})

	// TitleBar Actions
	t.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Raise()
	 	f.IsMousePressed = true
	 	f.MousePos = e.GlobalPos()
		f.Pos = f.Widget.Window().Pos()
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
		f.Widget.Window().Move(newPos)
	})

	t.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
		if f.IconMaximize.IsVisible() {
			f.windowMaximize()
		} else {
			f.windowRestore()
		}
	})

	// Button Actions
	f.IconMinimize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Window().SetWindowState(core.Qt__WindowMinimized)
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconMaximize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowMaximize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconRestore.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowRestore()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.IconClose.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
	})
}

func(f *QFramelessWindow) windowMaximize() {
	f.IconMaximize.SetVisible(false)
	f.IconRestore.SetVisible(true)
	f.Widget.Window().SetWindowState(core.Qt__WindowMaximized)
}

func(f *QFramelessWindow) windowRestore() {
	f.IconMaximize.SetVisible(true)
	f.IconRestore.SetVisible(false)
	f.Widget.Window().SetWindowState(core.Qt__WindowNoState)
}
