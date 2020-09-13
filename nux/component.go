// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

import "github.com/nuxui/nuxui/log"

type Component interface {
	Widget
	Component() Widget
	Content() Widget
}

type component struct {
	WidgetBase
	component Widget
	content   Widget
}

/*
NewComponent component, child not nil
*/
func NewComponent(compt, content Widget) Component {
	if compt == nil {
		log.Fatal("nux", "component can not ne nil")
	}

	if content == nil {
		log.Fatal("nux", "child can not ne nil")
	}

	return &component{
		component: compt,
		content:   content,
	}
}

func (me *component) Component() Widget {
	return me.component
}

func (me *component) Content() Widget {
	return me.content
}
