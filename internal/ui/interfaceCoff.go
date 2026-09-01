package ui

import (
	"context"
	"log/slog"
	"sort"
	"time"
	"voicer/internal/storage"
	"voicer/internal/ui/coffee"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
)

func (ui *UI) PointerEvents(gtx layout.Context) {
	for {
		keyEvent, ok := gtx.Event(
			key.Filter{
				Name: key.NameF3,
			},
			key.Filter{
				Name: key.NameF4,
			},
			key.Filter{
				Name: key.NameF5,
			},
			key.Filter{
				Name: key.NameF6,
			},
			key.Filter{
				Name:     "F",
				Optional: key.ModCtrl,
			},
			key.Filter{
				Name:     "E",
				Optional: key.ModCtrl,
			},
		)
		if !ok {
			break
		}

		switch kE := keyEvent.(type) {
		case key.Event:
			if kE.State == key.Press {
				switch kE.Name {
				case key.NameF3:
					if !ui.a.DevelopMode.Load() {
						ui.a.DevelopMode.Store(true)

						ui.coffeeMenu.HighBlockMobile = 0.7
						ui.coffeeMenu.HighBlockSettingMobile = 0.3

						ui.SetValuesIntoEditors()

						ctx, cancel := context.WithCancel(context.Background())
						ui.a.AutoSaveCancel = cancel

						go ui.a.AutoSave(ctx, 10*time.Second, ui.UpdateState)
						continue
					}
					ui.coffeeMenu.HighBlockMobile = 1.0
					ui.coffeeMenu.HighBlockSettingMobile = 0.0

					ui.a.AutoSaveCancelHelper()
					ui.a.DevelopMode.Store(false)
				case key.NameF4:
					if !ui.a.IsFullScreen.Load() {
						ui.a.IsFullScreen.Store(true)
						ui.a.Window.Option(app.Fullscreen.Option())
					} else {
						ui.a.IsFullScreen.Store(false)
						ui.a.Window.Option(app.Size(unit.Dp(425), unit.Dp(750)))

					}
				case key.NameF5:
					if !ui.a.Rendering.Load() {
						ui.rendPDF(gtx)
					}
				case key.NameF6:
					if !ui.a.Rendering.Load() {
						ui.rendPNG(gtx)
					}
				case "F":
					if kE.Modifiers.Contain(key.ModCtrl) {
						if ui.a.MobileMode.Load() {
							ui.a.MobileMode.Store(false)
							ui.SetValuesIntoEditors()
							continue
						}

						ui.a.MobileMode.Store(true)
						ui.SetValuesIntoEditors()
					}
				case "E":
					if kE.Modifiers.Contain(key.ModCtrl) {
						if !ui.a.DevelopModeEnh.Load() {
							ui.a.DevelopModeEnh.Store(true)
							continue
						}
						ui.a.DevelopModeEnh.Store(false)
					}
				}
			}
		}
	}
}

func (ui *UI) SetState(s storage.State) {
	ui.coffeeMenu.SetState(s.StateCoffee)
}

func (ui *UI) SetFlags() {
	ui.coffeeMenu.SetFlags()
}

func (ui *UI) UpdateState() storage.State {
	var state storage.State
	state.StateCoffee = ui.coffeeMenu.UpdateState()
	return state
}

func (ui *UI) SetValuesIntoEditors() {
	m := map[string]map[string]*storage.Type{}
	if ui.a.MobileMode.Load() {
		m = ui.coffeeMenu.Settings.FontSzM
	} else {
		m = ui.coffeeMenu.Settings.FontSz
	}

	ui.coffeeMenu.Settings.SettingItems = make([]coffee.SettingItem, len(m))
	var idx int

	for k1, v1 := range m {
		item := &ui.coffeeMenu.Settings.SettingItems[idx]
		item.Title = k1

		var idxParam int
		item.Parameters = make([]coffee.Parameter, len(v1))
		for k2, v2 := range v1 {
			param := &item.Parameters[idxParam]
			param.Title = k2
			param.Editor = new(widget.Editor)
			param.Editor.SetText(v2.Str)

			idxParam++
		}

		idx++
	}
	sort.Slice(ui.coffeeMenu.Settings.SettingItems, func(i, j int) bool {
		return ui.coffeeMenu.Settings.SettingItems[i].Title > ui.coffeeMenu.Settings.SettingItems[j].Title
	})

}

func (ui *UI) rendPDF(gtx layout.Context) {
	ui.a.Rendering.Store(true)
	ui.a.TargetRender = "PDF "

	ui.a.MobileMode.Store(true)
	ui.SetValuesIntoEditors()

	macro := op.Record(gtx.Ops)
	d := ui.coffeeMenu.RendMobile(gtx, ui.a.Th, ui.coffeeMenu.Settings.FontSzM)
	macro.Stop()

	go func() {
		err := ui.Screenshots(1037, 1843, d.Size.Y, 1)
		if err != nil {
			slog.Info(err.Error())
		}
		ui.a.MobileMode.Store(false)
		ui.SetValuesIntoEditors()
		ui.a.Rendering.Store(false)
		ui.coffeeMenu.ResultSavePdf(gtx, ui.a.TargetRender, err)
	}()
}

func (ui *UI) rendPNG(gtx layout.Context) {
	ui.a.Rendering.Store(true)
	ui.a.TargetRender = "PNG "

	ui.a.MobileMode.Store(false)
	ui.SetValuesIntoEditors()

	gtx.Constraints.Max.X = gtx.Dp(3840)

	macro := op.Record(gtx.Ops)
	d := ui.coffeeMenu.LengthContent(gtx, ui.a.Th, ui.coffeeMenu.Settings.FontSz)
	macro.Stop()

	go func() {
		err := ui.Screenshots(3840, 2160, d.Size.Y, 2)
		if err != nil {
			slog.Info(err.Error())
		}
		ui.a.Rendering.Store(false)
		ui.coffeeMenu.ResultSavePdf(gtx, ui.a.TargetRender, err)
	}()
}
