package coffee

import (
	"fmt"
	"image"
	"image/color"
	"time"
	"voicer/internal/logger"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

var widthDel = 400
var heightDel = 200
var widthLog = 7
var heightLog = 6

func (c *Components) RegisterModals(gtx C, th *material.Theme) {
	c.approveDeleteModal.Layout(gtx, th)
	c.loggingModal.Layout(gtx, th)
	c.savePdf.Layout(gtx, th)
}

func (c *Components) delModal(gtx C, s string) {
	c.approveDeleteModal.Widget = func(gtx C, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx C) layout.Dimensions {
			maxX := gtx.Dp(unit.Dp(widthDel))
			maxY := gtx.Dp(unit.Dp(heightDel))
			rect := image.Rect(0, 0, maxX, maxY)
			defer clip.RRect{Rect: rect, SE: gtx.Dp(10), SW: gtx.Dp(10), NW: gtx.Dp(10), NE: gtx.Dp(10)}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, colorWhite)

			gtx.Constraints.Max.X = maxX
			gtx.Constraints.Max.Y = maxY
			layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				space(0, 40),
				layout.Rigid(func(gtx C) layout.Dimensions {
					return material.Label(th, unit.Sp(18), "Удалить "+s+"?").Layout(gtx)
				}),
				space(0, 50),
				layout.Rigid(func(gtx C) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						space(30, 0),
						//confirm button
						layout.Flexed(0.5, func(gtx C) layout.Dimensions {
							return conditionalButton(gtx, th, c.Settings.confirmDeleting, "ok", unit.Dp(6))
						}),
						space(30, 0),
						//reject button
						layout.Flexed(0.5, func(gtx C) layout.Dimensions {
							return conditionalButton(gtx, th, c.Settings.rejectDeleting, "отмена", unit.Dp(6))
						}),
						space(30, 0),
					)
				}),
			)

			return layout.Dimensions{Size: rect.Max}
		})
	}
	c.approveDeleteModal.FinalAlpha = 245
	c.approveDeleteModal.Appear(gtx.Now)
}

func (c *Components) ResultSavePdf(gtx C, target string, err error) {
	c.savePdf.Widget = func(gtx C, th *material.Theme, anim *component.VisibilityAnimation) D {
		paint.FillShape(gtx.Ops, color.NRGBA{R: 241, G: 221, B: 188, A: 255}, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					if err == nil {
						return material.Label(th, unit.Sp(20), target+"сохранены в \"рабочий стол/images\"").Layout(gtx)
					}
					return material.Label(th, unit.Sp(20), "произошла ошибка при сохранении файлов").Layout(gtx)
				})
			}),
			space(0, 30),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = gtx.Dp(120)
				return material.Button(th, c.Settings.closeSavePDF, "закрыть").Layout(gtx)
			}),
		)
	}
	c.savePdf.FinalAlpha = 245
	c.savePdf.Appear(gtx.Now)
}

func (c *Components) logModal(gtx C, messages []logger.Message) {
	c.loggingModal.Widget = func(gtx C, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx C) layout.Dimensions {
			rect := image.Rect(0, 0, gtx.Constraints.Max.X-gtx.Constraints.Max.X/widthLog, gtx.Constraints.Max.Y-gtx.Constraints.Max.Y/heightLog)
			defer clip.RRect{Rect: rect, SE: gtx.Dp(10), SW: gtx.Dp(10), NW: gtx.Dp(10), NE: gtx.Dp(10)}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, colorWhite)

			c.listLog.Layout(gtx, len(messages), func(gtx C, i int) D {
				return logLine(gtx, th, messages[i])
			})
			gtx.Constraints.Max.X = gtx.Dp(unit.Dp(widthDel))
			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(heightDel))
			layout.Flex{Axis: layout.Vertical}.Layout(gtx)
			return layout.Dimensions{Size: rect.Max}
		})
	}
	c.loggingModal.FinalAlpha = 245
	c.loggingModal.Appear(gtx.Now)
}

func logLine(gtx C, th *material.Theme, message logger.Message) D {
	t, err := time.Parse(time.RFC3339Nano, message.Time)
	timeStr := message.Time
	if err == nil {
		timeStr = t.Format("02.01 15:04:05.000")
	}

	prefix := fmt.Sprintf("%s  %s line:%d ", timeStr, message.Source.Function, message.Source.Line)

	levelColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	msgColor := color.NRGBA{R: 255, G: 80, B: 80, A: 255}

	prefixLabel := material.Label(th, unit.Sp(12), prefix)
	prefixLabel.Color = levelColor
	prefixLabel.Font.Typeface = "SteticaBlack"
	prefixLabel.Alignment = text.Start

	msgLabel := material.Label(th, unit.Sp(12), message.Message)
	msgLabel.Color = msgColor
	msgLabel.Font.Typeface = "SteticaBlack"
	msgLabel.Alignment = text.Start

	return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx,
		func(gtx C) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) layout.Dimensions {
					return prefixLabel.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) layout.Dimensions {
					return msgLabel.Layout(gtx)
				}),
			)
		},
	)
}
