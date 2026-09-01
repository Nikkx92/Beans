package ui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"
	"voicer/internal/app"
	"voicer/internal/ui/coffee"

	"gioui.org/font/gofont"
	"gioui.org/gpu/headless"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type Screen int

const (
	CoffeeScreen Screen = iota
)

type UI struct {
	a             *app.App
	currentScreen Screen
	coffeeMenu    *coffee.Components
}

func NewUI(a *app.App) *UI {

	return &UI{
		a:             a,
		currentScreen: CoffeeScreen,
		coffeeMenu:    coffee.NewComponents(a),
	}

}

func loader(gtx layout.Context, s string) {
	tempTh := material.NewTheme()
	tempTh.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	paint.FillShape(gtx.Ops, color.NRGBA{R: 241, G: 221, B: 188, A: 255}, clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Op())
	layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Label(tempTh, 20, "Генерация "+s+" ...").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: 30}.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max = image.Pt(100, 100)
			return material.Loader(tempTh).Layout(gtx)
		}),
	)
}

func (ui *UI) Run(gtx layout.Context, th *material.Theme) {
	if ui.a.Rendering.Load() {
		loader(gtx, ui.a.TargetRender)
		return
	}

	ui.coffeeMenu.HandleEvents(gtx)

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			scaleDesc := float32(ui.a.ScreenWidth) / 1920.0
			gtx.Metric = unit.Metric{
				PxPerDp: scaleDesc,
				PxPerSp: scaleDesc,
			}

			switch ui.currentScreen {
			case CoffeeScreen:
				if ui.a.MobileMode.Load() {
					scaleMob := float32(ui.a.ScreenWidth) / 1037.0
					gtx.Metric = unit.Metric{
						PxPerDp: scaleMob,
						PxPerSp: scaleMob,
					}

					return ui.coffeeMenu.MobileView(gtx, th, ui.coffeeMenu.Settings.FontSzM)
				}

				return ui.coffeeMenu.DescScreen(gtx, th, ui.coffeeMenu.Settings.FontSz)
			}
			return layout.Dimensions{}
		}))

	ui.coffeeMenu.RegisterModals(gtx, th)
}

func (ui *UI) Screenshots(width, height, length int, scale float32) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("не удалось получить домашнюю директорию: %w", err)
	}

	tss := ui.a.TargetRender + time.Now().Format("2006.01.02 15.04.05")
	basePath := filepath.Join(homeDir, "Desktop", "images", tss) + string(filepath.Separator)
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку images: %w", err)
	}

	sz := image.Point{X: width, Y: height}
	w, err := headless.NewWindow(sz.X, sz.Y)
	if err != nil {
		return fmt.Errorf("не удалось создать headless-окно: %w", err)
	}
	defer w.Release()

	ops := new(op.Ops)
	pageCounter := 1

	img := image.NewRGBA(image.Rectangle{Max: sz})

	for currentScrollY := 0; currentScrollY < length; currentScrollY += height {
		ops.Reset()

		gtx := layout.Context{
			Ops: ops,
			Metric: unit.Metric{
				PxPerDp: scale,
				PxPerSp: scale,
			},
			Constraints: layout.Exact(sz),
		}

		stack := op.Offset(image.Pt(0, -currentScrollY)).Push(gtx.Ops)

		gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(length))

		switch ui.a.TargetRender {
		case "PDF ":
			ui.coffeeMenu.RendMobile(gtx, ui.a.Th, ui.coffeeMenu.Settings.FontSzM)
		case "PNG ":
			ui.coffeeMenu.RendDesk(gtx, ui.a.Th, ui.coffeeMenu.Settings.FontSz)
		}

		stack.Pop()

		w.Frame(gtx.Ops)

		if err := w.Screenshot(img); err != nil {
			return fmt.Errorf("ошибка скриншота на странице %d: %w", pageCounter, err)
		}

		var imgBuf bytes.Buffer
		if err := png.Encode(&imgBuf, img); err != nil {
			return fmt.Errorf("ошибка кодирования PNG на странице %d: %w", pageCounter, err)
		}

		switch ui.a.TargetRender {
		case "PDF ":
			currentPdfPath := fmt.Sprintf("%s_page_%d.pdf", basePath, pageCounter)
			outF, err := os.OpenFile(currentPdfPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}

			singleImgReader := []io.Reader{bytes.NewReader(imgBuf.Bytes())}
			err = api.ImportImages(nil, outF, singleImgReader, nil, nil)
			if err != nil {
				outF.Close()
				return err
			}
			outF.Close()
		case "PNG ":
			currentPngPath := fmt.Sprintf("%s_page_%d.png", basePath, pageCounter)
			if err := os.WriteFile(currentPngPath, imgBuf.Bytes(), 0o666); err != nil {
				fmt.Println("Ошибка записи файла:", err)
			}
		}

		pageCounter++
	}

	ops = nil
	runtime.GC()
	debug.FreeOSMemory()

	return nil
}
