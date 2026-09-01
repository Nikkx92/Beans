package coffee

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var hFontTitleSetting = 26
var wSettingsEditors = 50

func (c *Components) drawSettings(gtx C, th *material.Theme, item int) D {
	settingItem := &c.Settings.SettingItems[item]
	title := settingItem.Title

	return layout.Inset{Top: 15, Bottom: 15}.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return material.Label(th, 16, title).Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					renderParameters(th, settingItem.Parameters)...,
				)
			}),
		)

	})
}

func renderParameters(th *material.Theme, parameters []Parameter) []layout.FlexChild {
	children := make([]layout.FlexChild, len(parameters))

	for i := range parameters {
		v := &parameters[i]
		itm := v.Title

		children[i] = layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Label(th, 16, itm).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					ed := material.Editor(th, v.Editor, "16")

					ed.Editor.Alignment = text.Middle
					ed.TextSize = unit.Sp(hFontTitleSetting)
					return widget.Border{
						Color:        colorBlack,
						CornerRadius: 6,
						Width:        1,
					}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(wSettingsEditors))
						return ed.Layout(gtx)

					})
				}),
				space(0, 50),
			)
		})
	}
	return children
}
