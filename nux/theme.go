// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nux

type Theme interface {
	GetAttr(widgetName, themeName, themeKind, styleName string) Attr
}

var appTheme Theme

func UseTheme(theme Theme) {
	appTheme = theme
}
