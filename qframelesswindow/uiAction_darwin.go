package qframelesswindow

import (
        "github.com/therecipe/qt/core"
        "github.com/therecipe/qt/gui"
)

func (f *QFramelessWindow) SetTitleBarActions() {
	t := f.TitleBar

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
		if f.BtnMaximize.IsVisible() {
			f.windowMaximize()
		} else {
			f.windowRestore()
		}
	})

	// Button Actions
	f.BtnMinimize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.Widget.Window().SetWindowState(core.Qt__WindowMinimized)
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnMaximize.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowMaximize()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnRestore.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		f.windowRestore()
		f.Widget.Hide()
		f.Widget.Show()
	})

	f.BtnClose.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
	})
}

func(f *QFramelessWindow) windowMaximize() {
	f.BtnMaximize.SetVisible(false)
	f.BtnRestore.SetVisible(true)
	f.Widget.Window().SetWindowState(core.Qt__WindowMaximized)
}

func(f *QFramelessWindow) windowRestore() {
	f.BtnMaximize.SetVisible(true)
	f.BtnRestore.SetVisible(false)
	f.Widget.Window().SetWindowState(core.Qt__WindowNoState)
}
