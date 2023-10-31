// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && !android) || (windows && cairo)

package nux

import (
	// "github.com/millken/nuxui/log"
	"runtime"

	"github.com/millken/nuxui/nux/internal/cairo"
	"github.com/millken/nuxui/nux/internal/pango"
)

type nativeFont struct {
	fd *pango.FontDescription
}

func createNativeFont(family string, traits uint32, weight FontWeight, size int32) *nativeFont {
	me := &nativeFont{
		fd: pango.FontDescriptionNew(),
	}
	me.fd.SetFamily(family)
	me.fd.SetSize(size * pango.Scale)
	me.fd.SetWeight(pango.WEIGHT_NORMAL) //TODO::
	runtime.SetFinalizer(me, freeNativeFont)
	return me
}

func freeNativeFont(me *nativeFont) {
	me.fd.Free()
}

func (me *nativeFont) SetFamily(family string) {
	me.fd.SetFamily(family)
}

func (me *nativeFont) SetSize(size int32) {
	me.fd.SetSize(size * pango.Scale)
}

func (me *nativeFont) SetWeight(weight int32) {
	me.fd.SetWeight(pango.WEIGHT_NORMAL) //TODO::
}

type nativeFontLayout struct {
	layout *pango.Layout
}

var emptyCairo = cairo.Create(nil)

func newNativeFontLayout() *nativeFontLayout {
	me := &nativeFontLayout{
		layout: pango.CairoCreateLayout(emptyCairo),
	}
	me.layout.SetWrap(pango.WRAP_WORD)
	runtime.SetFinalizer(me, freeNativeFontLayout)
	return me
}

func freeNativeFontLayout(me *nativeFontLayout) {
	me.layout.Free()
}

func (me *nativeFontLayout) MeasureText(font Font, paint Paint, text string, width, height int32) (textWidth, textHeight int32) {
	// pango.TestLayout()
	me.layout.SetFontDescription(font.native().fd)
	me.layout.SetText(text)
	me.layout.SetWidth(width * pango.Scale)
	me.layout.SetHeight(height * pango.Scale)
	textWidth, textHeight = me.layout.GetPixelSize()
	return
}

func (me *nativeFontLayout) DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32) {
	cr := canvas.native().cairo
	r, g, b, a := paint.Color().RGBAf()
	cr.SetSourceRGBA(r, g, b, a)
	me.layout.SetFontDescription(font.native().fd)
	me.layout.SetText(text)
	me.layout.SetWidth(width * pango.Scale)
	me.layout.SetHeight(height * pango.Scale)
	me.layout.CairoShowLayout(cr)
}

func (me *nativeFontLayout) CharacterIndexForPoint(font Font, text string, width, height int32, x, y float32) uint32 {
	me.layout.SetFontDescription(font.native().fd)
	me.layout.SetText(text)
	me.layout.SetWidth(width * pango.Scale)
	me.layout.SetHeight(height * pango.Scale)
	// index, trailing, hit := me.layout.XYtoIndex(int32(x), int32(y))
	// log.I("nuxui", "CharacterIndexForPoint x=%f,y=%f index=%d, trailing=%d, hit=%t", x, y, index, trailing, hit)
	index, _, _ := me.layout.XYtoIndex(int32(x), int32(y))
	return uint32(index)
}
