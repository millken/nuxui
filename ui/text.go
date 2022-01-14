// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

// TODO Text can automatically fine-tune the spacing to ensure that the font occupies the entire line. Basic Text does not do this and uses the new AlignedText

import (
	"math"
	"time"

	"github.com/nuxui/nuxui/nux"
	"github.com/nuxui/nuxui/util"
)

type Text interface {
	nux.Widget
	nux.Size
	Visual

	Text() string
	SetText(text string)
}

type text struct {
	*nux.WidgetBase
	*nux.WidgetSize
	*WidgetVisual

	text               string
	textSize           float32
	textColor          nux.Color
	textHighlightColor nux.Color
	paint              nux.Paint
	ellipsize          int

	downTime time.Time
}

func NewText(context nux.Context, attrs ...nux.Attr) Text {
	attr := getAttr(attrs...)
	me := &text{
		text:               attr.GetString("text", ""),
		textSize:           attr.GetFloat32("textSize", 12),
		textColor:          attr.GetColor("textColor", nux.White),
		textHighlightColor: attr.GetColor("textHighlightColor", nux.Transparent),
		paint:              nux.NewPaint(attr.GetAttr("font", nux.Attr{})),
		// ellipsize: ellipsizeFromName(attr.GetString("ellipsize", "none")),
	}

	me.WidgetBase = nux.NewWidgetBase(context, me, attrs...)
	me.WidgetSize = nux.NewWidgetSize(context, me, attrs...)
	me.WidgetVisual = NewWidgetVisual(context, me, attrs...)
	me.WidgetSize.AddOnSizeChanged(me.onSizeChanged)
	me.WidgetVisual.AddOnVisualChanged(me.onVisualChanged)

	return me
}

func (me *text) OnMount(content nux.Widget) {
	nux.OnTapDown(me, me.onTapDown)
	nux.OnTapUp(me, me.onTapUp)
	nux.OnTapCancel(me, me.onTapUp)
}

func (me *text) onTapDown(detail nux.GestureDetail) {
	me.SetBackgroundColor(0xFF938276)
	me.downTime = time.Now()
	nux.NewTimerBackToUI(nux.GESTURE_DOWN2UP_DELAY*time.Millisecond, func() {

	})
}

func (me *text) onTapUp(detail nux.GestureDetail) {
	if sub := time.Since(me.downTime); sub < nux.GESTURE_DOWN2UP_DELAY*time.Millisecond {
		nux.NewTimerBackToUI(nux.GESTURE_DOWN2UP_DELAY*time.Millisecond-sub, func() {
			me.doTapUp(detail)
		})
	} else {
		me.doTapUp(detail)
	}
}

func (me *text) doTapUp(detail nux.GestureDetail) {
	me.SetBackgroundColor(0xFFFFFFFF)
}

func (me *text) onSizeChanged(widget nux.Widget) {

}
func (me *text) onVisualChanged(widget nux.Widget) {
	nux.RequestRedraw(me)
}

func (me *text) Text() string {
	return me.text
}

func (me *text) SetText(text string) {
	if me.text != text {
		me.text = text
		nux.RequestLayout(me)
	}
}

// Responsible for determining the position of the widget align, margin...
func (me *text) Layout(dx, dy, left, top, right, bottom int32) {
	// log.V("nuxui", "text layout %d, %d, %d, %d, %d, %d", dx, dy, left, top, right, bottom)
}

func (me *text) Measure(width, height int32) {
	// measuredDuration := log.Time()
	// defer log.TimeEnd(measuredDuration, "nuxui", "ui.Text Measure ")

	var vPPt float32 // horizontal padding percent
	var vPPx float32 // horizontal padding pixel
	var hPPt float32
	var hPPx float32

	ms := me.MeasuredSize()

	// 1. Calculate its own padding size
	if me.HasPadding() {
		switch me.PaddingLeft().Mode() {
		case nux.Pixel:
			l := me.PaddingLeft().Value()
			ms.Padding.Left = util.Roundi32(l)
			hPPx += l
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				l := me.PaddingLeft().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Left = util.Roundi32(l)
				hPPx += l
			case nux.Auto:
				hPPt += me.PaddingLeft().Value()
			}
		}

		switch me.PaddingRight().Mode() {
		case nux.Pixel:
			r := me.PaddingRight().Value()
			ms.Padding.Right = util.Roundi32(r)
			hPPx += r
		case nux.Percent:
			switch nux.MeasureSpecMode(width) {
			case nux.Pixel:
				r := me.PaddingRight().Value() / 100 * float32(nux.MeasureSpecValue(width))
				ms.Padding.Right = util.Roundi32(r)
				hPPx += r
			case nux.Auto:
				hPPt += me.PaddingRight().Value()
			}
		}

		switch me.PaddingTop().Mode() {
		case nux.Pixel:
			t := me.PaddingTop().Value()
			ms.Padding.Top = util.Roundi32(t)
			vPPx += t
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				t := me.PaddingTop().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Top = util.Roundi32(t)
				vPPx += t
			case nux.Auto:
				vPPt += me.PaddingTop().Value()
			}
		}

		switch me.PaddingBottom().Mode() {
		case nux.Pixel:
			b := me.PaddingBottom().Value()
			ms.Padding.Bottom = util.Roundi32(b)
			vPPx += b
		case nux.Percent:
			switch nux.MeasureSpecMode(height) {
			case nux.Pixel:
				b := me.PaddingBottom().Value() / 100 * float32(nux.MeasureSpecValue(height))
				ms.Padding.Bottom = util.Roundi32(b)
				vPPx += b
			case nux.Auto:
				vPPt += me.PaddingBottom().Value()
			}
		}
	}

	if nux.MeasureSpecMode(width) == nux.Auto || nux.MeasureSpecMode(height) == nux.Auto {
		w := int32(width)
		h := int32(height)

		me.paint.SetTextSize(me.textSize)
		outW, outH := me.paint.MeasureText(me.text, float32(nux.MeasureSpecValue(w)), float32(nux.MeasureSpecValue(h)))

		ms := me.MeasuredSize()
		if nux.MeasureSpecMode(width) == nux.Auto {
			w := (float32(outW) + hPPx) / (1.0 - hPPt/100.0)
			ms.Width = nux.MeasureSpec(int32(math.Ceil(float64(w))), nux.Pixel)
		} else {
			ms.Width = width
		}

		if nux.MeasureSpecMode(height) == nux.Auto {
			h := (float32(outH) + vPPx) / (1.0 - vPPt/100.0)
			ms.Height = nux.MeasureSpec(int32(math.Ceil(float64(h))), nux.Pixel)
		} else {
			ms.Height = height
		}
	}
}

func (me *text) Draw(canvas nux.Canvas) {
	if me.Background() != nil {
		me.Background().Draw(canvas)
	}

	ms := me.MeasuredSize()
	canvas.Save()
	canvas.Translate(float32(ms.Padding.Left), float32(ms.Padding.Top))
	canvas.ClipRect(0, 0,
		float32(ms.Width-ms.Padding.Left-ms.Padding.Right),
		float32(ms.Height-ms.Padding.Top-ms.Padding.Bottom))

	if me.text != "" {
		me.paint.SetTextSize(me.textSize)
		me.paint.SetColor(me.textColor)
		canvas.DrawText(me.text, float32(ms.Width), float32(ms.Height), me.paint)
	}

	canvas.Restore()

	if me.Foreground() != nil {
		me.Foreground().Draw(canvas)
	}
}
