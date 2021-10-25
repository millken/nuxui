// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"testing"

	"github.com/nuxui/nuxui/log"
	"github.com/nuxui/nuxui/nux"
)

var template = `
{
  import: {
    ui: github.com/nuxui/nuxui/ui,
  },

  layout: {
	id: "root",
	widget: ui.Column,
	width: 1wt,
	height: 1wt,
	background: #215896,
	padding: {left: 10px, top: 10px, right: 10px, bottom: 10px},
	children:[
	{
		id: "edit",
		widget: ui.Editor,
		width: 1wt,
		height: 30px,
		background: #982368,
		text: "nuxui.org example",
		font: {family: "Menlo, Monaco, Courier New, monospace", size: 14, color: #ffffff }
	},{
		id: "header",
		widget: ui.Column,
		width: 1wt,
		height: 100px,
		background: #123098,
		padding: {left: 10px, top: 10px, right: 10px, bottom: 10px},
		children:[
		{
			id: "xxx",
			widget: ui.Text,
			width: auto,
			height: auto,
			background: #982368,
			text: "{{me.name}}",
		}
		]
	}
	]
  }
}
  
  `

func TestAttr(t *testing.T) {
	defer log.Close()
	attr := nux.ParseAttr(template)
	log.V("test", "%s", attr)
}