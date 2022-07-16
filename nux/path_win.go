// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && !cairo

package nux

import (
	"nuxui.org/nuxui/nux/internal/win32"
	"runtime"
)

type path struct {
	ptr  *win32.GpPath
	curX float32
	curY float32
}

func newPath() *path {
	me := &path{}
	win32.GdipCreatePath(win32.FillModeAlternate, &me.ptr)
	runtime.SetFinalizer(me, freePath)
	return me
}

func freePath(me *path) {
	win32.GdipDeletePath(me.ptr)
}

func (me *path) native() *path {
	return me
}

func (me *path) Rect(x, y, width, height float32) {
	win32.GdipAddPathRectangle(me.ptr, x, y, width, height)
}

func (me *path) RoundRect(x, y, width, height, rx, ry float32) {
	dx := rx + rx
	dy := ry + ry

	win32.GdipAddPathLine(me.ptr, x+rx, y, x+width-dx, y)
	win32.GdipAddPathArc(me.ptr, x+width-dx, y, dx, dy, 270, 90)
	win32.GdipAddPathLine(me.ptr, x+width, y+ry, x+width, y+height-dy)
	win32.GdipAddPathArc(me.ptr, x+width-dx, y+height-dy, dx, dy, 0, 90)
	win32.GdipAddPathLine(me.ptr, x+width-dx, y+height, x+rx, y+height)
	win32.GdipAddPathArc(me.ptr, x, y+height-dy, dx, dy, 90, 90)
	win32.GdipAddPathLine(me.ptr, x, y+height-dy, x, y+ry)
	win32.GdipAddPathArc(me.ptr, x, y, dx, dy, 180, 90)
	win32.GdipClosePathFigure(me.ptr)
}

func (me *path) Ellipse(cx, cy, rx, ry float32) {
	win32.GdipAddPathEllipse(me.ptr, cx-rx, cy-ry, rx+rx, ry+ry)
}

func (me *path) MoveTo(x, y float32) {
	me.curX = x
	me.curY = y
}

func (me *path) LineTo(x, y float32) {
	win32.GdipAddPathLine(me.ptr, me.curX, me.curY, x, y)
	me.curX = x
	me.curY = y
}

func (me *path) CurveTo(x1, y1, x2, y2, x3, y3 float32) {
	win32.GdipAddPathBezier(me.ptr, me.curX, me.curY, x1, y1, x2, y2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) CurveToV(x2, y2, x3, y3 float32) {
	cx1 := me.curX + 2.0/3.0*(x2-me.curX)
	cy1 := me.curY + 2.0/3.0*(y2-me.curY)
	cx2 := x3 + 2.0/3.0*(x2-x3)
	cy2 := y3 + 2.0/3.0*(y2-y3)

	win32.GdipAddPathBezier(me.ptr, me.curX, me.curY, cx1, cy1, cx2, cy2, x3, y3)
	me.curX = x3
	me.curY = y3
}

func (me *path) Close() {
	win32.GdipClosePathFigure(me.ptr)
}
