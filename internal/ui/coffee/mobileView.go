package coffee

import (
	"image"
	"voicer/fonts"
	"voicer/internal/storage"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

var TotalLength int

func (c *Components) MobileView(gtx C, th *material.Theme, m map[string]map[string]*storage.Type) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		//blocks item
		layout.Flexed(c.HighBlockMobile, func(gtx C) D {
			paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())
			TotalLength = -9999999

			return c.ListBlocks.Layout(gtx, len(c.Blocks), func(gtx C, i int) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return logoMobile(gtx, c.icons["лого_гориз"])
					}),
					layout.Rigid(func(gtx C) D {
						return c.drawHeader(gtx, th, i, m)
					}),
				)
			})
		}),
		//settings
		layout.Flexed(c.HighBlockSettingMobile, func(gtx C) D {
			if c.a.DevelopMode.Load() {
				paint.FillShape(gtx.Ops, colorWhite, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())

				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						grid := component.Grid(th, c.Grid).Layout(gtx, 1, len(m),
							func(axis layout.Axis, index, constraint int) int {
								if axis == layout.Horizontal {
									return gtx.Dp(220)
								}
								return gtx.Dp(600)
							},
							func(gtx C, row, col int) D {
								return c.drawSettings(gtx, th, col)
							},
						)
						return grid
					}),
				)
			}
			return D{}
		}),
	)

}

func (c *Components) RendMobile(gtx C, th *material.Theme, m map[string]map[string]*storage.Type) D {
	gtx.Constraints.Max.X = gtx.Dp(1037)
	gtx.Constraints.Min.X = gtx.Dp(1037)
	gtx.Constraints.Max.Y = gtx.Dp(1000000)

	gtx.Metric.PxPerDp = 1.0
	gtx.Metric.PxPerSp = 1.0

	TotalLength = 0

	paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		func() []layout.FlexChild {
			children := make([]layout.FlexChild, len(c.Blocks))
			for i := 0; i < len(c.Blocks); i++ {
				children[i] = layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						//fill last area
						layout.Rigid(func(gtx C) D {
							macro := op.Record(gtx.Ops)
							firstSortSizeY := c.highSortMobile(gtx, th, i, 0, m).Size.Y
							macro.Stop()

							headerPx := gtx.Dp(unit.Dp(m["h_header"]["value"].Num))

							if TotalLength+headerPx+firstSortSizeY > gtx.Dp(1843) {
								return fillLastArea(gtx, 1843)
							}

							return D{}
						}),
						//logo
						layout.Rigid(func(gtx C) D {
							if TotalLength == 0 || i == 0 {
								return logoMobile(gtx, c.icons["лого_гориз"])
							}
							return D{}
						}),
						//header
						layout.Rigid(func(gtx C) D {
							return c.drawHeader(gtx, th, i, m)
						}),
					)
				})
			}
			return children
		}()...,
	)
}

func (c *Components) highSortMobile(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(1, func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldSort, "сорт", fonts.SteticaBlack, colorBlack, m)
										}),
										layout.Rigid(func(gtx C) D {
											return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldBlend, "купаж", fonts.SteticaItalic, color14011082, m)
										}),
									)
								}),
								layout.Flexed(float32(m["w_column3"]["value"].Num)/float32(100), func(gtx C) D {
									if c.Blocks[blockIdx].Sorts[sortIdx].fieldBlend.Text() == "Арабика" {
										return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												c.Blocks[blockIdx].Sorts[sortIdx].fieldGrade.Alignment = text.Alignment(layout.End)
												return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldGrade, "грейд", fonts.SteticaBlack, color14011082, m)
											}),
											layout.Rigid(func(gtx C) D {
												return c.icon(gtx, "грейд", m)
											}),
										)
									}
									return D{}
								}),
							)
						}),
						//sort
						layout.Rigid(func(gtx C) D {
							HighColumns := func(gtx C) D {
								macro := op.Record(gtx.Ops)
								gtx.Constraints.Min = image.Point{}

								columnsDims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Flexed(float32(m["w_column1"]["value"].Num)/float32(100), func(gtx C) D {
										return c.column2(gtx, th, blockIdx, sortIdx, m)
									}),
									layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
									layout.Flexed(float32(m["w_column2"]["value"].Num)/float32(100), func(gtx C) D {
										return c.column3(gtx, th, blockIdx, sortIdx, m)
									}),
									layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
									layout.Flexed(float32(m["w_column3"]["value"].Num)/float32(100), func(gtx C) D {
										return c.column4(gtx, th, blockIdx, sortIdx, m)
									}),
								)
								columnsCall := macro.Stop()

								return layout.Stack{}.Layout(gtx,
									//columns
									layout.Stacked(func(gtx C) D {
										columnsCall.Add(gtx.Ops)
										return columnsDims
									}),
									//vertical lines
									layout.Stacked(func(gtx C) D {
										w1 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column1"]["value"].Num) / float32(100)))
										w2 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column2"]["value"].Num) / float32(100)))
										lineWidth := gtx.Dp(WidthVerticalLine)

										op.Offset(image.Pt(w1, 0)).Add(gtx.Ops)
										lineVertical(gtx, 0, 0, columnsDims.Size.Y)

										op.Offset(image.Pt(w2+lineWidth, 0)).Add(gtx.Ops)
										lineVertical(gtx, 0, 0, columnsDims.Size.Y)

										return D{Size: columnsDims.Size}
									}),
								)
							}(gtx)
							return HighColumns
						}),
						//fill 30px after sort
						layout.Rigid(func(gtx C) D {
							paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 30)}.Op())
							return D{Size: image.Pt(gtx.Constraints.Max.X, 30)}
						}),
					)
				}),
				//limited edition
				layout.Stacked(func(gtx C) D {
					if c.Blocks[blockIdx].Sorts[sortIdx].limited.Value {
						return c.icon(gtx, "limited", m)
					}
					return D{}
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Spacer{Width: unit.Dp(wButtonsSection)}.Layout(gtx)
		}),
	)
}

func (c *Components) DrawSortsMobile(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	macro := op.Record(gtx.Ops)
	totalHighSortMobile := c.highSortMobile(gtx, th, blockIdx, sortIdx, m)
	sortCall := macro.Stop()

	isFirstSort := sortIdx == 0

	var lineHigh int
	if !isFirstSort {
		lineHigh = gtx.Dp(HighHorizontalLine)
	}

	willPageBreak := TotalLength+totalHighSortMobile.Size.Y+lineHigh > gtx.Dp(1843)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		//next page
		layout.Rigid(func(gtx C) D {
			if willPageBreak {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					//fill
					layout.Rigid(func(gtx C) D {
						return fillLastArea(gtx, 1843)
					}),
					//logo
					layout.Rigid(func(gtx C) D {
						return logoMobile(gtx, c.icons["лого_гориз"])
					}),
					//head
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, func(gtx C) D {
								return c.head(gtx, th, blockIdx, m)
							}),
							space(wButtonsSection, 0),
						)
					}),
				)
			}

			return D{}
		}),
		//horizontal line
		layout.Rigid(func(gtx C) D {
			if !willPageBreak && !isFirstSort {
				gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Dp(unit.Dp(wButtonsSection))
				dims := lineHorizontal(gtx)
				TotalLength += dims.Size.Y
				return dims
			}
			return D{}
		}),
		//sort
		layout.Rigid(func(gtx C) D {
			TotalLength += gtx.Dp(unit.Dp(totalHighSortMobile.Size.Y))
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				//sorts
				layout.Flexed(1, func(gtx C) D {
					sortCall.Add(gtx.Ops)
					return totalHighSortMobile
				}),
			)
		}),
		//fill last space
		layout.Rigid(func(gtx C) D {
			if blockIdx == len(c.Blocks)-1 && sortIdx == len(c.Blocks[blockIdx].Sorts)-1 {
				paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(1843-TotalLength)))}.Op())
				return D{Size: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(1843-TotalLength)))}
			}
			return D{}
		}),
	)

}

func logoMobile(gtx C, logo *widget.Image) D {
	logo.Fit = widget.Contain

	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, colorBlack)
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Inset{
				Top: unit.Dp(40), Bottom: unit.Dp(40),
				Left: unit.Dp(40), Right: unit.Dp(40),
			}.Layout(gtx, func(gtx C) D {
				return logo.Layout(gtx)
			})
		}),
	)
	TotalLength += dims.Size.Y
	return dims
}

func fillLastArea(gtx C, length int) D {
	// Высота страницы в пикселях
	pageHeight := gtx.Dp(unit.Dp(length)) // 2160 пикселей при scale 2.0

	// Сколько места осталось на текущей странице
	remainingSpace := pageHeight - TotalLength
	paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, remainingSpace)}.Op())
	TotalLength = 0

	return D{Size: image.Pt(gtx.Constraints.Max.X, remainingSpace)}
}
