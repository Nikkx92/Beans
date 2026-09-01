package coffee

import (
	"image"
	"voicer/fonts"
	"voicer/internal/storage"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var wButtonsSection = 60
var upArrow = "↑"
var downArrow = "↓"
var squares = [][]int{
	{1, 0},
	{3, 0},
	{4, 0},
	{5, 0},
	{6, 0},
	{8, 0},
	{10, 0},
	{2, 1},
	{3, 1},
	{4, 1},
	{6, 1},
	{8, 1},
	{9, 1},
	{1, 2},
	{2, 2},
	{3, 2},
	{5, 2},
	{7, 2},
	{10, 2},
	{0, 3},
	{1, 3},
	{3, 3},
	{0, 4},
	{2, 4},
	{6, 4},
	{0, 5},
}

func (c *Components) head(gtx C, th *material.Theme, idx int, m map[string]map[string]*storage.Type) D {
	TotalLength += gtx.Dp(unit.Dp(m["h_header"]["value"].Num))

	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		layout.Expanded(func(gtx C) D {

			rect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Dp(unit.Dp(m["h_header"]["value"].Num)))

			defer clip.Rect(rect).Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, color14011082)

			return D{Size: rect.Max}
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				//flag
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: 60, Right: 15}.Layout(gtx, func(gtx C) D {
						if c.Blocks[idx].loader {
							gtx.Constraints = layout.Exact(image.Pt(gtx.Dp(30), gtx.Dp(30)))
							return material.Loader(th).Layout(gtx)
						}

						if c.Blocks[idx].flag == nil {
							return conditionalButton(gtx, th, c.Blocks[idx].seekFlag, "seek flag", 8)
						}

						return c.flag(gtx, idx, m)

					})
				}),
				//title
				layout.Rigid(func(gtx C) D {
					return c.editorField(gtx, th, c.Blocks[idx].fieldCountry, "страна", fonts.SteticaBlack, colorWhite, m)
				}),
				//pixels
				layout.Rigid(func(gtx C) D {

					var pixels []layout.FlexChild
					for _, v := range squares {
						pixels = append(pixels, layout.Rigid(func(gtx C) D {

							return square(gtx, v[0], v[1], m["cell_size"]["value"].Num, m["h_header"]["value"].Num)
						}))
					}
					return layout.Flex{}.Layout(gtx, pixels...)
				}),
			)
		}),
	)
}

func (c *Components) drawHeader(gtx C, th *material.Theme, blockIdx int, m map[string]map[string]*storage.Type) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		//block head
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				//header
				layout.Flexed(1, func(gtx C) D {
					return c.head(gtx, th, blockIdx, m)
				}),
				//buttons
				layout.Rigid(func(gtx C) D {
					if c.a.DevelopMode.Load() && !c.a.MobileMode.Load() {
						return c.headButtons(gtx, th, blockIdx)
					}

					return layout.Spacer{Width: unit.Dp(wButtonsSection)}.Layout(gtx)
				}),
			)

		}),
		//sorts
		layout.Rigid(func(gtx C) D {
			sorts := make([]layout.FlexChild, len(c.Blocks[blockIdx].Sorts))
			for i := range c.Blocks[blockIdx].Sorts {
				sorts[i] = layout.Rigid(func(gtx C) D {
					if c.a.MobileMode.Load() {
						return c.DrawSortsMobile(gtx, th, blockIdx, i, m)
					}

					return c.drawSorts(gtx, th, blockIdx, i, m)
				})
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, sorts...)
		}),
	)
}

func (c *Components) headButtons(gtx C, th *material.Theme, blockIdx int) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		//spacer
		space(1, 0),
		//up/down buttons
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = 70
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return conditionalButton(gtx, th, c.Blocks[blockIdx].upCountryBtn, upArrow, 0)
				}),
				space(0, 1),
				layout.Rigid(func(gtx C) D {
					return conditionalButton(gtx, th, c.Blocks[blockIdx].downCountryBtn, downArrow, 0)
				}),
			)
		}),
		//spacer
		space(1, 0),
		//deleteCountry/addSort buttons
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = 140
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return conditionalButton(gtx, th, c.Blocks[blockIdx].deleteCountry, "удалить страну", 0)
				}),
				space(0, 1),
				layout.Rigid(func(gtx C) D {
					return conditionalButton(gtx, th, c.Blocks[blockIdx].addSort, "добавить сорт", 0)
				}),
			)
		}),
	)
}

func (c *Components) column1(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		//sort
		layout.Rigid(func(gtx C) D {
			return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldSort, "сорт", fonts.SteticaBlack, colorBlack, m)
		}),
		//blend
		layout.Rigid(func(gtx C) D {
			return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldBlend, "купаж", fonts.SteticaItalic, color14011082, m)
		}),
		//grade
		layout.Rigid(func(gtx C) D {
			if c.Blocks[blockIdx].Sorts[sortIdx].fieldBlend.Text() == "Арабика" {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					//grade img
					layout.Rigid(func(gtx C) D {
						return c.icon(gtx, "грейд", m)
					}),
					//grade
					layout.Rigid(func(gtx C) D {
						c.Blocks[blockIdx].Sorts[sortIdx].fieldGrade.Alignment = text.Alignment(layout.Start)
						return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldGrade, "грейд", fonts.SteticaBlack, color14011082, m)
					}),
					layout.Rigid(func(gtx C) D {
						if c.Blocks[blockIdx].Sorts[sortIdx].limited.Value {
							return c.icon(gtx, "limited", m)
						}
						return D{}
					}),
				)
			} else if c.Blocks[blockIdx].Sorts[sortIdx].fieldBlend.Text() == "Робуста" {
				if c.Blocks[blockIdx].Sorts[sortIdx].limited.Value {
					return layout.Inset{Left: unit.Dp(175)}.Layout(gtx, func(gtx C) D {
						return c.icon(gtx, "limited", m)
					})
				}
			}
			return D{}
		}),
	)
}

func (c *Components) column2(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		//location
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return c.icon(gtx, "локация", m)
				}),
				layout.Rigid(func(gtx C) D {
					return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldLocation, "локация", fonts.SteticaMedium, color14011082, m)
				}),
			)
		}),
		//above
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return c.icon(gtx, "высота", m)
				}),
				layout.Rigid(func(gtx C) D {
					return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldHeightAbove, "высота", fonts.SteticaRegular, colorBlack, m)
				}),
			)
		}),
		//drying
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: 5, Bottom: 5}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return c.icon(gtx, "сушка", m)
					}),
					layout.Rigid(func(gtx C) D {
						return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldDrying, "сушка", fonts.SteticaRegular, colorBlack, m)
					}),
				)
			})
		}),
		//cornSize
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return c.icon(gtx, "размер", m)
				}),
				layout.Rigid(func(gtx C) D {
					return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldCornSize, "размер", fonts.SteticaRegular, colorBlack, m)
				}),
			)
		}),
	)
}

func (c *Components) column3(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		//taste
		layout.Rigid(func(gtx C) D {
			return c.icon(gtx, "вкус", m)
		}),
		layout.Rigid(func(gtx C) D {
			return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldTaste, "вкус", fonts.SteticaRegular, colorBlack, m)
		}),
	)
}

func (c *Components) column4(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	if c.Blocks[blockIdx].Sorts[sortIdx].inStock.Value {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			//price
			layout.Rigid(func(gtx C) D {
				return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldPrice, "цена", fonts.SteticaBlack, colorBlack, m)
			}),
			//quantity
			layout.Rigid(func(gtx C) D {
				return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldQuantity, "кол-во", fonts.SteticaItalic, color14011082, m)
			}),
		)
	}
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return c.editorField(gtx, th, c.Blocks[blockIdx].Sorts[sortIdx].fieldStockDesc, "наличие", fonts.SteticaBoldItalic, color14011082, m)
		}))
}
