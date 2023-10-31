// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && !android

package nux

import (
	"runtime"

	"github.com/millken/nuxui/nux/internal/cairo"
	"github.com/millken/nuxui/nux/internal/linux/xlib"
)

type canvas struct {
	cairo   *cairo.Cairo
	surface *cairo.Surface
}

func canvasFromWindow(display *xlib.Display, drawable xlib.Drawable, visual *xlib.Visual, width, height int32) *canvas {
	surface := cairo.XlibSurfaceCreate(display, drawable, visual, width, height)
	me := &canvas{
		cairo:   cairo.Create(surface),
		surface: surface,
	}
	runtime.SetFinalizer(me, freeCanvas)
	return me
}

func freeCanvas(me *canvas) {
	me.cairo.Destroy()
	me.surface.Destroy()
}
