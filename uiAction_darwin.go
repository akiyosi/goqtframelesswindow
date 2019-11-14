package qframelesswindow

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

func (f *QFramelessWindow) SetupTitleBarActions() {
	t := f.TitleBar

	// TitleBar Actions
	t.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Raise()
		f.IsTitleBarPressed = true
		f.TitleBarMousePos = e.GlobalPos()
		f.Position = f.Pos()
	})

	t.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	t.ConnectMouseMoveEvent(func(e *gui.QMouseEvent) {
		if !f.IsTitleBarPressed {
			return
		}
		x := f.Position.X() + e.GlobalPos().X() - f.TitleBarMousePos.X()
		y := f.Position.Y() + e.GlobalPos().Y() - f.TitleBarMousePos.Y()
		newPos := core.NewQPoint2(x, y)
		f.Move(newPos)
	})

	t.ConnectMouseDoubleClickEvent(func(e *gui.QMouseEvent) {
		if f.BtnMaximize.IsVisible() {
			f.WindowMaximize()
		} else {
			f.WindowRestore()
		}
	})

	// Button Actions
	f.BtnMinimize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnMaximize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnRestore.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnClose.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.IsTitleBarPressed = false
	})

	f.BtnMinimize.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.WindowMinimize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnMaximize.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.WindowMaximize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnRestore.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.WindowRestore()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnClose.ConnectMouseReleaseEvent(func(e *gui.QMouseEvent) {
		f.Close()
	})
}

func (f *QFramelessWindow) WindowMinimize() {
	f.SetWindowState(core.Qt__WindowMinimized)
}

func (f *QFramelessWindow) WindowMaximize() {
	f.BtnMaximize.SetVisible(false)
	f.BtnRestore.SetVisible(true)
	f.Layout.SetContentsMargins(0, 0, 0, 0)
	f.SetWindowState(core.Qt__WindowMaximized)
}

func (f *QFramelessWindow) WindowRestore() {
	f.BtnMaximize.SetVisible(true)
	f.BtnRestore.SetVisible(false)
	f.Layout.SetContentsMargins(f.shadowMargin, f.shadowMargin, f.shadowMargin, f.shadowMargin)
	f.SetWindowState(core.Qt__WindowNoState)
}
