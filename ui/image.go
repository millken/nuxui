// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Image interface {
	nux.Widget
	nux.Size
	nux.Layout
	nux.Measure
	nux.Draw
	Visual

	Src() string
	SetSrc(src string)
	ScaleType() ScaleType
	SetScaleType(scaleType ScaleType)
}

type ScaleType int32

const (
	ScaleType_Matrix ScaleType = iota
	ScaleType_Center
	ScaleType_CenterCrop
	ScaleType_CenterInside
	ScaleType_FitXY
	ScaleType_FitStart
	ScaleType_FitCenter
	ScaleType_FitEnd
)

type Repeat int32

func NewImage(attrs ...nux.Attr) Image {
	attr := nux.MergeAttrs(attrs...)
	me := &image{
		scaleX:    1.0,
		scaleY:    1.0,
		offsetX:   0,
		offsetY:   0,
		scaleType: convertScaleTypeFromString(attr.GetString("scaleType", "matrix")),
		src:       attr.GetString("src", ""),
	}
	me.WidgetBase = nux.NewWidgetBase(attrs...)
	me.WidgetSize = nux.NewWidgetSize(attrs...)
	me.WidgetVisual = NewWidgetVisual(me, attrs...)
	me.WidgetSize.AddSizeObserver(me.onSizeChanged)
	if me.src != "" {
		me.srcDrawable = NewImageDrawableWithSource(me.src)
	}
	return me
}

type image struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	scaleType   ScaleType
	src         string
	srcDrawable ImageDrawable
	scaleX      float32
	scaleY      float32
	offsetX     float32
	offsetY     float32
}

func convertScaleTypeFromString(scaleType string) ScaleType {
	switch scaleType {
	case "matrix":
		return ScaleType_Matrix
	case "center":
		return ScaleType_Center
	case "centerCrop":
		return ScaleType_CenterCrop
	case "centerInside":
		return ScaleType_CenterInside
	case "fitXY":
		return ScaleType_FitXY
	case "fitStart":
		return ScaleType_FitStart
	case "fitCenter":
		return ScaleType_FitCenter
	case "fitEnd":
		return ScaleType_FitEnd
	}

	log.Fatal("nux", fmt.Sprintf("unknow scale type %s, only support 'matrix', 'center', 'centerCrop', 'centerInside', 'fitXY', 'fitStart', 'fitEnd'", scaleType))
	return ScaleType_Center
}

func (me *image) OnMount() {
	if me.src != "" {
		me.srcDrawable = NewImageDrawableWithSource(me.src)
	}
}

// TODO if not have Layout, then use default layout to set frame
func (me *image) Layout(x, y, width, height int32) {
	frame := me.Frame()

	var imgW, imgH float32
	if me.srcDrawable != nil {
		w, h := me.srcDrawable.Size()
		imgW = float32(w)
		imgH = float32(h)
	}
	innerW := frame.Width - frame.Padding.Left - frame.Padding.Right
	innerH := frame.Height - frame.Padding.Top - frame.Padding.Bottom

	if imgW == 0 || imgH == 0 || innerW == 0 || innerH == 0 {
		me.scaleX = 1.0
		me.scaleY = 1.0
		me.offsetX = 0
		me.offsetY = 0
		return
	}

	switch me.scaleType {
	case ScaleType_Matrix:
		me.scaleX = 1.0
		me.scaleY = 1.0
		me.offsetX = 0
		me.offsetY = 0
	case ScaleType_Center:
		me.scaleX = 1.0
		me.scaleY = 1.0
		me.offsetX = (float32(innerW) - imgW) / 2
		me.offsetY = (float32(innerH) - imgH) / 2
	case ScaleType_CenterCrop:
		r := imgW / imgH
		ir := float32(innerW) / float32(innerH)
		if ir > r {
			newH := float32(innerW) / r
			me.scaleX = float32(innerW) / imgW
			me.scaleY = newH / imgH
			me.offsetX = 0
			me.offsetY = (float32(innerH) - newH) / 2
		} else {
			newW := float32(innerH) * r
			me.scaleX = newW / imgW
			me.scaleY = float32(innerH) / imgH
			me.offsetX = (float32(innerW) - newW) / 2
			me.offsetY = 0
		}
	case ScaleType_CenterInside:
		if imgW > float32(innerW) || imgH > float32(innerH) {
			r := imgW / imgH
			r2 := float32(innerW) / float32(innerH)
			if r2 > r {
				newW := float32(innerH) * r
				me.scaleX = newW / imgW
				me.scaleY = float32(innerH) / imgH
				me.offsetX = (float32(innerW) - newW) / 2
				me.offsetY = 0
			} else {
				newH := float32(innerW) / r
				me.scaleX = float32(innerW) / imgW
				me.scaleY = newH / imgH
				me.offsetX = 0
				me.offsetY = (float32(innerH) - newH) / 2
			}
		} else {
			me.scaleX = 1.0
			me.scaleY = 1.0
			me.offsetX = (float32(innerW) - imgW) / 2
			me.offsetY = (float32(innerH) - imgH) / 2
		}

	case ScaleType_FitXY:
		me.scaleX = float32(innerW) / imgW
		me.scaleY = float32(innerH) / imgH
		me.offsetX = 0
		me.offsetY = 0
	case ScaleType_FitCenter, ScaleType_FitStart, ScaleType_FitEnd:
		r := imgW / imgH
		r2 := float32(innerW) / float32(innerH)
		if r2 > r {
			newW := float32(innerH) * r
			me.scaleX = newW / imgW
			me.scaleY = float32(innerH) / imgH

			switch me.scaleType {
			case ScaleType_FitStart:
				me.offsetX = 0
				me.offsetY = 0
			case ScaleType_FitCenter:
				me.offsetX = (float32(innerW) - newW) / 2
				me.offsetY = 0
			case ScaleType_FitEnd:
				me.offsetX = float32(innerW) - newW
				me.offsetY = 0
			}
		} else {
			newH := float32(innerW) / r
			me.scaleX = float32(innerW) / imgW
			me.scaleY = newH / imgH

			switch me.scaleType {
			case ScaleType_FitStart:
				me.offsetX = 0
				me.offsetY = 0
			case ScaleType_FitCenter:
				me.offsetX = 0
				me.offsetY = (float32(innerH) - newH) / 2
			case ScaleType_FitEnd:
				me.offsetX = 0
				me.offsetY = float32(innerH) - newH
			}
		}
	}
}

func (me *image) Measure(width, height int32) {
	if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {
		frame := me.Frame()
		dw, dh := me.srcDrawable.Size()
		if nux.MeasureSpecMode(width) == nux.Auto && me.srcDrawable != nil {
			frame.Width = nux.MeasureSpec(dw+frame.Padding.Left+frame.Padding.Right, nux.Pixel)
		} else {
			frame.Width = width
		}

		if nux.MeasureSpecMode(height) == nux.Auto && me.srcDrawable != nil {
			frame.Height = nux.MeasureSpec(dh+frame.Padding.Top+frame.Padding.Bottom, nux.Pixel)
		} else {
			frame.Height = height
		}
	}

	if me.Padding() != nil {
		frame := me.Frame()

		switch me.Padding().Left.Mode() {
		case nux.Pixel:
			frame.Padding.Left = util.Roundi32(me.Padding().Left.Value())
		}

		switch me.Padding().Top.Mode() {
		case nux.Pixel:
			frame.Padding.Top = util.Roundi32(me.Padding().Top.Value())
		}

		switch me.Padding().Right.Mode() {
		case nux.Pixel:
			frame.Padding.Right = util.Roundi32(me.Padding().Right.Value())
		}

		switch me.Padding().Bottom.Mode() {
		case nux.Pixel:
			frame.Padding.Bottom = util.Roundi32(me.Padding().Bottom.Value())
		}
	}
}

func (me *image) onSizeChanged() {
	nux.RequestLayout(me)
}

func (me *image) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	frame := me.Frame()
	canvas.Save()
	canvas.Translate(float32(frame.X), float32(frame.Y))
	canvas.Translate(float32(frame.Padding.Left), float32(frame.Padding.Top))

	if me.srcDrawable != nil {
		canvas.Translate(me.offsetX, me.offsetY)
		canvas.Scale(me.scaleX, me.scaleY)
		me.srcDrawable.Draw(canvas)
	}
	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}

func (me *image) Src() string {
	return me.src
}

func (me *image) SetSrc(src string) {
	if me.src == src {
		return
	}

	me.src = src

	if me.src != "" {
		me.srcDrawable = NewImageDrawableWithSource(me.src)
	}

	nux.RequestLayout(me)
	nux.RequestRedraw(me)
}

func (me *image) ScaleType() ScaleType {
	return me.scaleType
}

func (me *image) SetScaleType(scaleType ScaleType) {
	me.scaleType = scaleType
}
