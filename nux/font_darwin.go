// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin && !ios

package nux

import (
	"math"
	"runtime"

	"github.com/millken/nuxui/nux/internal/darwin"
)

type nativeFont struct {
	ptr darwin.NSFont
}

func createNativeFont(family string, traits uint32, weight FontWeight, size int32) *nativeFont {
	font := darwin.SharedNSFontManager().FontWithFamily(family, uint32(traits), int32(weight), float32(size))
	if font.IsNil() {
		font = darwin.NSFontSystemFontOfSize(float32(size), fontWeightToNative(weight))
	}
	me := &nativeFont{
		ptr: font,
	}
	runtime.SetFinalizer(me, freeNativeFont)
	return me
}

func freeNativeFont(me *nativeFont) {
	darwin.NSObject_release(uintptr(me.ptr))
}

func (me *nativeFont) SetFamily(family string) {
	// me.fd.SetFamily(family)
}

func (me *nativeFont) SetSize(size int32) {
	// me.fd.SetSize(size * pango.Scale)
}

func (me *nativeFont) SetWeight(weight int32) {
	// me.fd.SetWeight(pango.WEIGHT_NORMAL) //TODO::
}

type nativeFontLayout struct {
	layout    darwin.NSLayoutManager
	container darwin.NSTextContainer
}

func newNativeFontLayout() *nativeFontLayout {
	me := &nativeFontLayout{
		layout:    darwin.NewNSLayoutManager(),
		container: darwin.NewNSTextContainer(0, 0),
	}
	me.layout.AddTextContainer(me.container)
	runtime.SetFinalizer(me, freeNativeFontLayout)
	return me
}

func freeNativeFontLayout(me *nativeFontLayout) {
	me.layout.RemoveTextContainerAtIndex(0)
	darwin.NSObject_release(uintptr(me.container))
	darwin.NSObject_release(uintptr(me.layout))
}

func (me *nativeFontLayout) MeasureText(font Font, paint Paint, text string, width, height int32) (textWidth, textHeight int32) {
	w, h := me.layout.MeasureText(me.container, font.native().ptr, text, float32(width), float32(height))
	return int32(math.Ceil(float64(w))), int32(math.Ceil(float64(h)))
}

func (me *nativeFontLayout) DrawText(canvas Canvas, font Font, paint Paint, text string, width, height int32) {
	me.layout.DrawText(me.container, font.native().ptr, text, float32(width), float32(height), uint32(paint.Color()), 0) //TODO:: bgcolor
}

func (me *nativeFontLayout) CharacterIndexForPoint(font Font, text string, width, height int32, x, y float32) uint32 {
	index, fraction := me.layout.CharacterIndexForPoint(me.container, font.native().ptr, text, float32(width), float32(height), x, y)
	if fraction > 0.5 {
		index++
	}
	return index
}

func fontWeightToNative(weight FontWeight) darwin.NSFontWeight {
	switch weight {
	case FontWeight_Thin:
		return darwin.NSFontWeightThin
	case FontWeight_ExtraLight:
		return darwin.NSFontWeightUltraLight
	case FontWeight_Light:
		return darwin.NSFontWeightLight
	case FontWeight_Normal:
		return darwin.NSFontWeightRegular
	case FontWeight_Medium:
		return darwin.NSFontWeightMedium
	case FontWeight_SemiBold:
		return darwin.NSFontWeightSemibold
	case FontWeight_Bold:
		return darwin.NSFontWeightBold
	case FontWeight_ExtraBold:
		return darwin.NSFontWeightHeavy
	case FontWeight_Black:
		return darwin.NSFontWeightBlack
	}
	return darwin.NSFontWeightRegular
}
