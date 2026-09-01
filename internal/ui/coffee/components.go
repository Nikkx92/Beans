package coffee

import (
	"image"
	"image/color"
	"strconv"
	"voicer/images"
	a "voicer/internal/app"
	"voicer/internal/storage"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type ItemBlock struct {
	flag           *widget.Image
	loader         bool
	seekFlag       *widget.Clickable
	fieldCountry   *widget.Editor
	addSort        *widget.Clickable
	deleteCountry  *widget.Clickable
	upCountryBtn   *widget.Clickable
	downCountryBtn *widget.Clickable
	Sorts          []SortBlock
}

type SortBlock struct {
	fieldSort        *widget.Editor
	fieldBlend       *widget.Editor
	fieldGrade       *widget.Editor
	fieldLocation    *widget.Editor
	fieldHeightAbove *widget.Editor
	fieldDrying      *widget.Editor
	fieldCornSize    *widget.Editor
	fieldTaste       *widget.Editor
	fieldPrice       *widget.Editor
	fieldQuantity    *widget.Editor
	fieldStockDesc   *widget.Editor
	deleteSort       *widget.Clickable
	up               *widget.Clickable
	down             *widget.Clickable
	inStock          *widget.Bool
	limited          *widget.Bool
}

var color14011082 = color.NRGBA{R: 140, G: 110, B: 82, A: 255}
var colorBlack = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
var colorWhite = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
var colorBackground = color.NRGBA{R: 241, G: 221, B: 188, A: 255}
var HighHorizontalLine = unit.Dp(2)
var WidthVerticalLine = unit.Dp(1.5)

type Components struct {
	a                      *a.App
	deleteFunc             func()
	Blocks                 []ItemBlock
	ListBlocks             *widget.List
	ListSettings           *widget.List
	listLog                *widget.List
	Grid                   *component.GridState
	approveDeleteModal     *component.ModalLayer
	loggingModal           *component.ModalLayer
	savePdf                *component.ModalLayer
	flags                  map[string]*widget.Image
	icons                  map[string]*widget.Image
	HighBlockMobile        float32
	HighBlockSettingMobile float32
	isAnimating            bool
	animatingFunc          func()
	Settings               SettingsComponents
}

type SettingsComponents struct {
	addCountry      *widget.Clickable
	closeSettings   *widget.Clickable
	confirmDeleting *widget.Clickable
	rejectDeleting  *widget.Clickable
	logBtn          *widget.Clickable
	closeSavePDF    *widget.Clickable
	SettingItems    []SettingItem
	FontSz          map[string]map[string]*storage.Type
	FontSzM         map[string]map[string]*storage.Type
}

func NewComponents(a *a.App) *Components {
	icons := images.LoadIcons(images.IconsFS)

	FontSz, FontSzM := NewSettingItems()

	flags := images.LoadFlags(images.FlagsFS)
	images.LoadAiCreatedFlags(flags)

	return &Components{
		a: a,
		ListBlocks: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		ListSettings: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		listLog: &widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		Grid:                   &component.GridState{},
		Blocks:                 make([]ItemBlock, len(a.State.StateCoffee.Items)),
		approveDeleteModal:     component.NewModal(),
		loggingModal:           component.NewModal(),
		savePdf:                component.NewModal(),
		icons:                  icons,
		flags:                  flags,
		HighBlockMobile:        1.0,
		HighBlockSettingMobile: 0.0,
		Settings: SettingsComponents{
			addCountry:      new(widget.Clickable),
			closeSettings:   new(widget.Clickable),
			confirmDeleting: new(widget.Clickable),
			rejectDeleting:  new(widget.Clickable),
			logBtn:          new(widget.Clickable),
			closeSavePDF:    new(widget.Clickable),
			FontSz:          FontSz,
			FontSzM:         FontSzM,
		},
	}
}

func (c *Components) editorField(gtx C, th *material.Theme, field *widget.Editor, hint string, font font.Typeface, col color.NRGBA, m map[string]map[string]*storage.Type) D {
	return layout.Inset{Top: unit.Dp(m[hint]["i_top_txt"].Num), Bottom: unit.Dp(m[hint]["i_down_txt"].Num), Left: unit.Dp(m[hint]["i_left_txt"].Num), Right: unit.Dp(m[hint]["i_right_txt"].Num)}.Layout(gtx, func(gtx C) D {
		localTh := *th
		localTh.Face = font

		effectiveHint := hint
		if field.Text() != "" {
			effectiveHint = ""
		}

		e := material.Editor(&localTh, field, effectiveHint)
		e.TextSize = unit.Sp(m[hint]["kegel"].Num)
		e.Color = col
		e.Editor.ReadOnly = !c.a.DevelopMode.Load()
		return e.Layout(gtx)
	})
}

func (c *Components) icon(gtx C, s string, m map[string]map[string]*storage.Type) D {
	return layout.Inset{Top: unit.Dp(m[s]["i_top_img"].Num), Bottom: unit.Dp(m[s]["i_down_img"].Num), Left: unit.Dp(m[s]["i_left_img"].Num), Right: unit.Dp(m[s]["i_right_img"].Num)}.Layout(gtx, func(gtx C) D {
		gtx.Constraints = layout.Exact(image.Pt(gtx.Dp(unit.Dp(m[s]["w_img"].Num)), gtx.Dp(unit.Dp(m[s]["h_img"].Num))))
		return c.icons[s].Layout(gtx)
	})
}

func (c *Components) flag(gtx layout.Context, blockIdx int, m map[string]map[string]*storage.Type) layout.Dimensions {

	gtx.Constraints = layout.Exact(image.Pt(gtx.Dp(unit.Dp(m["страна"]["w_img"].Num)), gtx.Dp(unit.Dp(m["страна"]["h_img"].Num))))

	return c.Blocks[blockIdx].seekFlag.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: unit.Dp(m["страна"]["i_top_img"].Num), Bottom: unit.Dp(m["страна"]["i_down_img"].Num), Left: unit.Dp(m["страна"]["i_left_img"].Num), Right: unit.Dp(m["страна"]["i_right_img"].Num)}.Layout(gtx, func(gtx C) D {
			return c.Blocks[blockIdx].flag.Layout(gtx)
		})
	})
}

func conditionalButton(gtx C, th *material.Theme, btn *widget.Clickable, text string, radius unit.Dp) D {
	b := material.Button(th, btn, text)
	b.CornerRadius = radius
	return b.Layout(gtx)
}

func lineHorizontal(gtx C) D {
	h := gtx.Dp(HighHorizontalLine)

	padd := gtx.Constraints.Max.X

	defer clip.Rect{Max: image.Pt(padd, h)}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color14011082}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{Size: image.Pt(padd, h)}
}

func lineVertical(gtx C, up, down int, height int) D {
	h := gtx.Dp(WidthVerticalLine)
	upDp := gtx.Dp(unit.Dp(up))
	downDp := gtx.Dp(unit.Dp(down))

	start := upDp
	finish := height - downDp

	if finish <= start {
		return D{}
	}

	defer clip.Rect{
		Min: image.Pt(0, start),
		Max: image.Pt(h, finish),
	}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color14011082}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return D{Size: image.Pt(h, height)}
}

func square(gtx C, row, col, cellSize, high int) D {
	cellSizeDp := gtx.Dp(unit.Dp(cellSize))
	highDp := gtx.Dp(unit.Dp(high))

	dimX0 := gtx.Constraints.Max.X - row*cellSizeDp
	dimX1 := dimX0 - cellSizeDp
	dimY1 := highDp - col*cellSizeDp
	dimY0 := dimY1 - cellSizeDp

	paint.FillShape(gtx.Ops, colorBackground, clip.Rect{
		Min: image.Pt(dimX1, dimY0),
		Max: image.Pt(dimX0, dimY1),
	}.Op())

	return D{Size: image.Pt(cellSizeDp, highDp)}
}

func space(w, h int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Spacer{Width: unit.Dp(w), Height: unit.Dp(h)}.Layout(gtx)
	})
}

func aToI(s string) int {
	d, err := strconv.Atoi(s)
	if err != nil {
		d = 26
	}
	return d
}
