// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/nuxui/nuxui/nux"
)

type Button interface {
	Label
}

type button label

func NewButton(attrs ...nux.Attr) Button {
	attr := nux.Attr{
		"selectable": false,
	}
	a := nux.MergeAttrs(attr, nux.MergeAttrs(attrs...))
	me := NewLabel(a)
	return Button(me)
}
