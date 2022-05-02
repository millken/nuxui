// Copyright 2018 The NuxUI Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package theme

import (
	"nuxui.org/nuxui/nux"
)

type Material struct {
}

func (me *Material) GetAttr(widgetName, themeName, themeKind, styleName string) nux.Attr {
	switch styleName {
	case "button":
		return nux.Attr{
			// "widget":     ui.Column,
			"width":      nux.ADimen(0, nux.Auto),
			"height":     nux.ADimen(0, nux.Auto),
			"textSize":   14,
			"textColor":  "#de000000",
			"textShadow": nux.Attr{"color": "#88000000", "x": 0, "y": 2, "blur": 3},
			"padding":    nux.Attr{"left": 16, "top": 8, "right": 16, "bottom": 8},
			"background": nux.Attr{
				"drawable": "nuxui.org/nuxui/ui.ShapeDrawable",
				"states": []nux.Attr{
					{
						"state": "default",
						"shape": nux.Attr{
							"shape":        "rect",
							"solid": "#ffe0e0e0",
							"cornerRadius": 4,
							"shadow":       nux.Attr{"color": "#88000000", "x": 0, "y": 1, "blur": 3},
						},
					},
					{
						"state": "pressed",
						"shape": nux.Attr{
							"shape":        "rect",
							"solid":        "#ff9e9e9e",
							"cornerRadius": 4,
							"shadow":       nux.Attr{"color": "#88000000", "x": 0, "y": 1, "blur": 3},
						},
					},
				},
			},
		}
	}
	return nil
}
