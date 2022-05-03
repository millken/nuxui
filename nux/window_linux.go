// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <cairo/cairo.h>
#include <cairo/cairo-xlib.h>

#include <X11/Xlib.h>

void window_getSize(Display* display, Window window, int32_t *width, int32_t *height);
void window_getContentSize(Display* display, Window window, int32_t *width, int32_t *height);
void window_setTitle(Display* display, Window window, char *name);
*/
import "C"

import (
	"unsafe"

	"nuxui.org/nuxui/log"
)

type window struct {
	windptr C.Window
	display *C.Display
	visual  *C.Visual
	surface *C.cairo_surface_t

	decor        Parent
	buffer       unsafe.Pointer
	bufferWidth  int32
	bufferHeight int32
	delegate     WindowDelegate
	focusWidget  Widget

	initEvent PointerEvent
	timer     Timer

	width  int32
	height int32
	title  string
}

func newWindow(attr Attr) *window {
	me := &window{}

	me.CreateDecor(attr)
	return me
}

func (me *window) CreateDecor(attr Attr) Widget {
	creator := FindRegistedWidgetCreator("nuxui.org/nuxui/ui.Layer")
	w := creator(attr)
	if p, ok := w.(Parent); ok {
		me.decor = p
	} else {
		log.Fatal("nuxui", "decor must is a Parent")
	}

	decorWindowList[w] = me

	return me.decor
}

// func (me *window) Created() {
// 	main := App().Manifest().Main()
// 	if main == "" {
// 		log.Fatal("nuxui", "no main widget found.")
// 	} else {
// 		mainWidgetCreator := FindRegistedWidgetCreator(main)
// 		ctx := &context{}
// 		// widgetTree := RenderWidget(mainWidgetCreator(ctx, Attr{}))
// 		widgetTree := mainWidgetCreator(ctx, Attr{})
// 		me.decor.AddChild(widgetTree)
// 	}

// 	me.excuteCreated(me.decor)
// }

// func (me *window) excuteCreated(widget Widget) {
// 	if c, ok := widget.(Created); ok {
// 		c.Created()
// 	}

// 	if c, ok := widget.(Component); ok {
// 		me.excuteCreated(c.Content())
// 	}

// 	if p, ok := widget.(Parent); ok {
// 		for _, child := range p.Children() {
// 			me.excuteCreated(child)
// 		}
// 	}
// }

// func (me *window) CreateDecor(attr Attr) Widget {
// 	creator := FindRegistedWidgetCreator("nuxui.org/nuxui/ui.Layer")
// 	w := creator(ctx, attr)
// 	if p, ok := w.(Parent); ok {
// 		me.decor = p
// 	} else {
// 		log.Fatal("nuxui", "decor must is a Parent")
// 	}

// 	decorWindowList[w] = me

// 	return me.decor
// }

// func (me *window) Measure(width, height MeasureDimen) {
// 	if me.decor == nil {
// 		return
// 	}

// 	me.width = width
// 	me.height = height

// 	if s, ok := me.decor.(Size); ok {
// 		if s.Frame().Width == width && s.Frame().Height == height {
// 			// return
// 		}

// 		s.Frame().Width = width
// 		s.Frame().Height = height
// 	}

// 	me.surfaceResized = true

// 	if f, ok := me.decor.(Measure); ok {
// 		f.Measure(width, height)
// 	}
// }

// func (me *window) Layout(dx, dy, left, top, right, bottom int32) {
// 	if me.decor == nil {
// 		return
// 	}

// 	if f, ok := me.decor.(Layout); ok {
// 		f.Layout(dx, dy, left, top, right, bottom)
// 	}
// }
var TestDraw func(Canvas)

func (me *window) Draw(canvas Canvas) {
	if me.decor != nil {
		if f, ok := me.decor.(Draw); ok {
			canvas.Save()
			if TestDraw != nil {
				TestDraw(canvas)
			} else {
				f.Draw(canvas)
			}
			canvas.Restore()
		}
	}
}

func (me *window) ID() uint64 {
	return 0
}

func (me *window) Size() (width, height int32) {
	var w, h C.int32_t
	C.window_getSize(me.display, me.windptr, &w, &h)
	return int32(w), int32(h)
}

func (me *window) ContentSize() (width, height int32) {
	var w, h C.int32_t
	C.window_getContentSize(me.display, me.windptr, &w, &h)
	return int32(w), int32(h)
}

func (me *window) LockCanvas() Canvas {
	w, h := me.ContentSize()
	me.surface = C.cairo_xlib_surface_create(me.display, me.windptr, me.visual, C.int(w), C.int(h))
	return newCanvas(me.surface)
}

func (me *window) UnlockCanvas(c Canvas) {
	// me.surface.GetCanvas().Destroy()
	// me.surface.Flush()
	// me.surface.Destroy()
	// C.XFlush(me.display)
}

func (me *window) Decor() Widget {
	return me.decor
}

// TODO:: use int? set 0.5 and return 127,0.498039
func (me *window) Alpha() float32 {
	// TODO::
	return 1.0
}

func (me *window) SetAlpha(alpha float32) {
	// TODO::
}

func (me *window) Title() string {
	return me.title
}

func (me *window) SetTitle(title string) {
	me.title = title
	ctitle := C.CString(title)
	C.window_setTitle(me.display, me.windptr, ctitle)
	C.free(unsafe.Pointer(ctitle))
}

func (me *window) SetDelegate(delegate WindowDelegate) {
	me.delegate = delegate
}

func (me *window) Delegate() WindowDelegate {
	return me.delegate
}

func (me *window) handlePointerEvent(e PointerEvent) {
	me.switchFocusIfPossible(e)

	if me.delegate != nil {
		if f, ok := me.delegate.(windowDelegate_HandlePointerEvent); ok {
			f.HandlePointerEvent(e)
			return
		}
	}

	hitTestResultManagerInstance.handlePointerEvent(me.Decor(), e)
}

func (me *window) handleScrollEvent(e ScrollEvent) {
	hitTestResultManagerInstance.handleScrollEvent(me.Decor(), e)
}

func (me *window) handleKeyEvent(e KeyEvent) {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(KeyEventHandler); ok {
			if f.OnKeyEvent(e) {
				return
			} else {
				goto other
			}
		}
	} else {
		goto other
	}

other:
	if me.decor != nil {
		me.handleOtherWidgetKeyEvent(me.decor, e)
	}
}

func (me *window) handleOtherWidgetKeyEvent(p Parent, e KeyEvent) bool {
	if p.ChildrenCount() > 0 {
		var compt Widget
		for _, c := range p.Children() {
			compt = nil
			if cpt, ok := c.(Component); ok {
				c = cpt.Content()
				compt = cpt
			}
			if cp, ok := c.(Parent); ok {
				if me.handleOtherWidgetKeyEvent(cp, e) {
					return true
				}
			} else if f, ok := c.(KeyEventHandler); ok {
				if f.OnKeyEvent(e) {
					return true
				}
			}

			if compt != nil {
				if f, ok := compt.(KeyEventHandler); ok {
					if f.OnKeyEvent(e) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (me *window) handleTypeEvent(e TypeEvent) {
	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(TypeEventHandler); ok {
			f.OnTypeEvent(e)
			return
		}
	}

	log.E("nuxui", "none widget handle typing event")
}

func (me *window) requestFocus(widget Widget) {
	if me.focusWidget == widget {
		return
	}

	if me.focusWidget != nil {
		if f, ok := me.focusWidget.(Focus); ok && f.HasFocus() {
			f.FocusChanged(false)
		}
	}

	if f, ok := widget.(Focus); ok {
		if f.Focusable() {
			me.focusWidget = widget
			if !f.HasFocus() {
				f.FocusChanged(true)
			}
		}
	}
}

func (me *window) switchFocusIfPossible(event PointerEvent) {
	if event.Type() != Type_PointerEvent || !event.IsPrimary() {
		return
	}

	switch event.Action() {
	case Action_Down:
		me.initEvent = event
		if event.Kind() == Kind_Mouse {
			me.switchFocusAtPoint(event.X(), event.Y())
		}
	case Action_Move:
		if event.Kind() == Kind_Touch {
			if me.timer != nil {
				if me.initEvent.Distance(event.X(), event.Y()) >= GESTURE_MIN_PAN_DISTANCE {
					me.switchFocusAtPoint(me.initEvent.X(), me.initEvent.Y())
				}
			}
		}
	case Action_Up:
		if event.Kind() == Kind_Touch {
			if event.Time().UnixNano()-me.initEvent.Time().UnixNano() < GESTURE_LONG_PRESS_TIMEOUT.Nanoseconds() {
				me.switchFocusAtPoint(event.X(), event.Y())
			}
		}
	}

}

func (me *window) switchFocusAtPoint(x, y float32) {
	if me.focusWidget != nil {
		if s, ok := me.focusWidget.(Size); ok {
			frame := s.Frame()
			if x >= float32(frame.X) && x <= float32(frame.X+frame.Width) &&
				y >= float32(frame.Y) && y <= float32(frame.Y+frame.Height) {
				// point is in current focus widget, do not need change focus
				return
			}
		}

		if f, ok := me.focusWidget.(Focus); ok && f.HasFocus() {
			me.focusWidget = nil
			f.FocusChanged(false)
		}
	}

}
