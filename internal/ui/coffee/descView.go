package coffee

import (
	"image"
	"voicer/internal/storage"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var wSettingsButtons = 170
var paddingTopHead = unit.Dp(20)

func (c *Components) DescScreen(gtx C, th *material.Theme, m map[string]map[string]*storage.Type) D {
	TotalLength = -9999999
	/*defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	r := image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
	area := clip.Rect(r).Push(gtx.Ops)
	event.Op(gtx.Ops, &c.a.MousePos)
	area.Pop()*/

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		//blocks
		layout.Flexed(0.87, func(gtx C) D {
			paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())
			if c.isAnimating {
				scrollSpeed := 25.0

				c.ListBlocks.Position.Offset += int(scrollSpeed)

				gtx.Execute(op.InvalidateCmd{})

				c.animatingFunc()
			}

			return c.ListBlocks.Layout(gtx, len(c.Blocks), func(gtx C, i int) D {
				return c.drawHeader(gtx, th, i, m)
			})

		}),
		//Settings_logo
		layout.Flexed(0.13, func(gtx C) D {
			if c.a.DevelopMode.Load() {
				paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					space(0, 15),
					//addCountry button
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Max.X = gtx.Dp(unit.Dp(gtx.Constraints.Max.X * 6 / 10))
						return conditionalButton(gtx, th, c.Settings.addCountry, "добавить страну", 6)
					}),
					space(0, 20),
					//close developerMode button
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Max.X = gtx.Dp(unit.Dp(wSettingsButtons))
						return conditionalButton(gtx, th, c.Settings.closeSettings, "закрыть редактор", 6)
					}),
					space(0, 15),
					//Settings
					layout.Rigid(func(gtx C) D {
						if c.a.DevelopModeEnh.Load() {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								//log button
								layout.Rigid(func(gtx C) D {
									gtx.Constraints.Max.X = gtx.Dp(unit.Dp(wSettingsButtons))
									return conditionalButton(gtx, th, c.Settings.logBtn, "log", 6)
								}),
								space(0, 15),
								//fields
								layout.Rigid(func(gtx C) D {
									return c.ListSettings.Layout(gtx, len(c.Settings.SettingItems), func(gtx C, i int) D {
										return c.drawSettings(gtx, th, i)
									})
								}),
							)
						}
						return D{}
					}),
				)
			}

			return logoDesc(gtx, c.icons["лого"])
		}),
	)
}

func (c *Components) drawSorts(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	macro := op.Record(gtx.Ops)
	highSort := c.highSortDesc(gtx, th, blockIdx, sortIdx, m)
	sortCall := macro.Stop()

	isFirstSort := sortIdx == 0

	var lineHigh int
	if !isFirstSort {
		lineHigh = gtx.Dp(HighHorizontalLine)
	}

	willPageBreak := TotalLength+highSort.Size.Y+lineHigh > gtx.Dp(1080)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		//next page
		layout.Rigid(func(gtx C) D {
			if willPageBreak {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					//fill area
					layout.Rigid(func(gtx C) D {
						return fillLastArea(gtx, 1080)
					}),
					//padding top
					layout.Rigid(func(gtx C) D {
						return padTopHead(gtx)
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
			TotalLength += highSort.Size.Y
			sortCall.Add(gtx.Ops)
			return highSort
		}),
	)
}

func (c *Components) highSortDesc(gtx C, th *material.Theme, blockIdx, sortIdx int, m map[string]map[string]*storage.Type) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		//columns
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					macro := op.Record(gtx.Ops)
					gtx.Constraints.Min = image.Point{}

					columnsDims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(float32(m["w_column1"]["value"].Num)/float32(100), func(gtx C) D {
							return c.column1(gtx, th, blockIdx, sortIdx, m)
						}),
						// Пропускаем место под Первую линию (просто пустой отступ нужной ширины)
						layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
						// Колонка 2
						layout.Flexed(float32(m["w_column2"]["value"].Num)/float32(100), func(gtx C) D {
							return c.column2(gtx, th, blockIdx, sortIdx, m)
						}),
						// Пропускаем место под Вторую линию
						layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
						// Колонка 3
						layout.Flexed(float32(m["w_column3"]["value"].Num)/float32(100), func(gtx C) D {
							return c.column3(gtx, th, blockIdx, sortIdx, m)
						}),
						// Пропускаем место под третью линию
						layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
						//колонка 4
						layout.Flexed(float32(m["w_column4"]["value"].Num)/float32(100), func(gtx C) D {
							return c.column4(gtx, th, blockIdx, sortIdx, m)
						}),
					)
					columnsCall := macro.Stop()

					return layout.Stack{}.Layout(gtx,
						// Слой 1: Сам текст колонок (воспроизводим запись)
						layout.Stacked(func(gtx C) D {
							columnsCall.Add(gtx.Ops)
							return columnsDims
						}),
						// Слой 2: Вертикальные линии
						layout.Stacked(func(gtx C) D {
							// Вычисляем точную X-координату для линий на основе ширины колонок
							w1 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column1"]["value"].Num) / float32(100)))
							w2 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column2"]["value"].Num) / float32(100)))
							w3 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column3"]["value"].Num) / float32(100)))
							lineWidth := gtx.Dp(WidthVerticalLine)

							// Рисуем первую линию
							op.Offset(image.Pt(w1, 0)).Add(gtx.Ops)
							lineVertical(gtx, 20, 0, columnsDims.Size.Y)

							// Рисуем вторую линию (сдвиг считается от начала координат всего блока)
							op.Offset(image.Pt(w2+lineWidth, 0)).Add(gtx.Ops)
							lineVertical(gtx, 20, 0, columnsDims.Size.Y)

							op.Offset(image.Pt(w3+lineWidth, 0)).Add(gtx.Ops)
							lineVertical(gtx, 20, 0, columnsDims.Size.Y)

							return D{Size: columnsDims.Size}
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if c.a.DevelopMode.Load() {
						return c.sortButtons(gtx, th, blockIdx, sortIdx)
					}
					return layout.Spacer{Width: unit.Dp(wButtonsSection)}.Layout(gtx)
				}),
			)
		}),
		//fill 20px
		layout.Rigid(func(gtx C) D {
			paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 20)}.Op())
			return D{Size: image.Pt(gtx.Constraints.Max.X, 20)}
		}),
	)
	/*macro := op.Record(gtx.Ops)
	gtx.Constraints.Min = image.Point{}

	columnsDims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(float32(m["w_column1"]["value"].Num)/float32(100), func(gtx C) D {
			return c.column1(gtx, th, blockIdx, sortIdx, m)
		}),
		// Пропускаем место под Первую линию (просто пустой отступ нужной ширины)
		layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
		// Колонка 2
		layout.Flexed(float32(m["w_column2"]["value"].Num)/float32(100), func(gtx C) D {
			return c.column2(gtx, th, blockIdx, sortIdx, m)
		}),
		// Пропускаем место под Вторую линию
		layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
		// Колонка 3
		layout.Flexed(float32(m["w_column3"]["value"].Num)/float32(100), func(gtx C) D {
			return c.column3(gtx, th, blockIdx, sortIdx, m)
		}),
		// Пропускаем место под третью линию
		layout.Rigid(func(gtx C) D { return D{Size: image.Pt(gtx.Dp(WidthVerticalLine), 0)} }),
		//колонка 4
		layout.Flexed(float32(m["w_column4"]["value"].Num)/float32(100), func(gtx C) D {
			return c.column4(gtx, th, blockIdx, sortIdx, m)
		}),
	)
	columnsCall := macro.Stop()

	return layout.Stack{}.Layout(gtx,
		// Слой 1: Сам текст колонок (воспроизводим запись)
		layout.Stacked(func(gtx C) D {
			columnsCall.Add(gtx.Ops)
			return columnsDims
		}),
		// Слой 2: Вертикальные линии
		layout.Stacked(func(gtx C) D {
			// Вычисляем точную X-координату для линий на основе ширины колонок
			w1 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column1"]["value"].Num) / float32(100)))
			w2 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column2"]["value"].Num) / float32(100)))
			w3 := int(float32(gtx.Constraints.Max.X) * (float32(m["w_column3"]["value"].Num) / float32(100)))
			lineWidth := gtx.Dp(WidthVerticalLine)

			// Рисуем первую линию
			op.Offset(image.Pt(w1, 0)).Add(gtx.Ops)
			lineVertical(gtx, 20, 0, columnsDims.Size.Y)

			// Рисуем вторую линию (сдвиг считается от начала координат всего блока)
			op.Offset(image.Pt(w2+lineWidth, 0)).Add(gtx.Ops)
			lineVertical(gtx, 20, 0, columnsDims.Size.Y)

			op.Offset(image.Pt(w3+lineWidth, 0)).Add(gtx.Ops)
			lineVertical(gtx, 20, 0, columnsDims.Size.Y)

			return D{Size: columnsDims.Size}
		}),
	)*/
}

func (c *Components) RendDesk(gtx C, th *material.Theme, m map[string]map[string]*storage.Type) D {
	gtx.Constraints.Max.X = gtx.Dp(1920)

	// 1. Рассчитываем ширину колонок на основе пропорций (87% и 13%)
	totalWidth := gtx.Constraints.Max.X
	leftWidth := int(float32(totalWidth) * 0.87)
	rightWidth := totalWidth - leftWidth

	// 2. Изолированно рендерим левую часть, чтобы узнать её точную высоту
	leftGtx := gtx
	leftGtx.Constraints.Min.X = leftWidth
	leftGtx.Constraints.Max.X = leftWidth
	leftGtx.Constraints.Min.Y = 0
	//leftGtx.Constraints.Max.Y = 1e6 // Даем свободу развернуться для снимка

	var leftMacro op.CallOp
	var leftDims layout.Dimensions

	macro := op.Record(leftGtx.Ops)
	leftDims = c.LengthContent(leftGtx, th, m)
	leftMacro = macro.Stop()

	// Выводим левую часть на холст снимка
	leftMacro.Add(gtx.Ops)

	// 3. Вычисляем количество логотипов на основе итоговой высоты левой части
	contentHeight := leftDims.Size.Y
	step := gtx.Dp(1080)

	count := contentHeight / step
	if contentHeight%step != 0 || count == 0 {
		count++
	}

	// 4. Отрисовываем логотипы последовательно в правой части
	for i := 0; i < count; i++ {
		// Считаем сдвиг: по X — в начало правой колонки, по Y — на текущий шаг (i * 2160)
		yOffset := i * step

		// Сохраняем состояние матрицы трансформаций и сдвигаем координаты
		// defer .Pop() вернет координаты обратно в конце текущей итерации цикла
		trans := op.Offset(image.Pt(leftWidth, yOffset)).Push(gtx.Ops)

		// Конфигурируем ограничения под один конкретный логотип
		rightGtx := gtx
		rightGtx.Constraints.Min.X = rightWidth
		rightGtx.Constraints.Max.X = rightWidth
		rightGtx.Constraints.Min.Y = step
		rightGtx.Constraints.Max.Y = step

		// Отрисовываем логотип
		logoDesc(rightGtx, c.icons["лого"])

		// Восстанавливаем стек операций для следующего шага цикла
		trans.Pop()
	}

	// Возвращаем итоговые размеры сгенерированного макета
	return layout.Dimensions{
		Size: image.Pt(totalWidth, contentHeight),
	}
}

func (c *Components) LengthContent(gtx C, th *material.Theme, m map[string]map[string]*storage.Type) D {
	TotalLength = 0
	gtx.Constraints.Max.Y = gtx.Dp(1000000)
	gtx.Metric = unit.Metric{
		PxPerDp: 2,
		PxPerSp: 2,
	}

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		func() []layout.FlexChild {
			children := make([]layout.FlexChild, len(c.Blocks))
			for i := 0; i < len(c.Blocks); i++ {
				children[i] = layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						//fill area
						layout.Rigid(func(gtx C) D {
							macro := op.Record(gtx.Ops)
							firstSortSizeY := c.highSortDesc(gtx, th, i, 0, m).Size.Y
							macro.Stop()

							headerPx := gtx.Dp(unit.Dp(m["h_header"]["value"].Num))

							if TotalLength+headerPx+firstSortSizeY > gtx.Dp(1080) {
								return fillLastArea(gtx, 1080)
							}

							return D{}
						}),
						//header
						layout.Rigid(func(gtx C) D {
							paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								//padding top
								layout.Rigid(func(gtx C) D {
									if i == 0 || TotalLength == 0 {
										return padTopHead(gtx)
									}
									return D{}
								}),
								//header
								layout.Rigid(func(gtx C) D {
									return c.drawHeader(gtx, th, i, m)
								}),
							)
						}),
					)
				})
			}
			return children
		}()...,
	)

	return dims
}

func logoDesc(gtx C, logo *widget.Image) D {
	paint.FillShape(gtx.Ops, colorBlack, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())
	logo.Fit = widget.Contain
	logo.Position = layout.Center

	pTop := gtx.Constraints.Max.Y / int(5.0*gtx.Metric.PxPerDp)
	pBottom := gtx.Constraints.Max.Y / int(20*gtx.Metric.PxPerDp)

	return layout.Inset{Top: unit.Dp(pTop), Bottom: unit.Dp(pBottom)}.Layout(gtx, func(gtx C) D {
		return logo.Layout(gtx)
	})
}

func padTopHead(gtx C) D {
	paint.FillShape(gtx.Ops, colorBackground, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(paddingTopHead))}.Op())
	TotalLength += gtx.Dp(paddingTopHead)
	return D{Size: image.Pt(gtx.Constraints.Max.X, gtx.Dp(paddingTopHead))}
}

func (c *Components) sortButtons(gtx C, th *material.Theme, blockIdx, sortIdx int) D {
	gtx.Constraints.Min.X = gtx.Dp(60)
	item := &c.Blocks[blockIdx].Sorts[sortIdx]
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		//up
		layout.Rigid(func(gtx C) D {
			return conditionalButton(gtx, th, item.up, upArrow, 0)
		}),
		space(0, 1),
		//delete
		layout.Rigid(func(gtx C) D {
			return conditionalButton(gtx, th, item.deleteSort, "удалить сорт", 0)
		}),
		space(0, 1),
		//down
		layout.Rigid(func(gtx C) D {
			return conditionalButton(gtx, th, item.down, downArrow, 0)
		}),
		space(0, 10),
		//inStock switch
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Label(th, unit.Sp(18), "наличие").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return material.Switch(th, item.inStock, "inStock").Layout(gtx)
				}),
			)
		}),
		space(0, 10),
		//LE switch
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return material.Label(th, unit.Sp(18), "LE").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return material.Switch(th, item.limited, "limited").Layout(gtx)
				}),
			)
		}),
	)
}
