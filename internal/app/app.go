package app

import (
	"context"
	"image/color"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
	"voicer/fonts"
	"voicer/internal/client"
	"voicer/internal/storage"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/fstanis/screenresolution"
)

type ViewUI interface {
	Run(gtx layout.Context, th *material.Theme)
	SetState(s storage.State)
	SetFlags()
	PointerEvents(gtx layout.Context)
}
type App struct {
	Window         *app.Window
	Th             *material.Theme
	IsFullScreen   atomic.Bool
	DevelopMode    atomic.Bool
	DevelopModeEnh atomic.Bool
	MobileMode     atomic.Bool
	Rendering      atomic.Bool
	State          storage.State
	AutoSaveCancel context.CancelFunc
	Client         *client.Client
	ScreenWidth    int
	TargetRender   string
}

func NewApp() *App {
	tr := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := &http.Client{Transport: tr, Timeout: 30 * time.Second}

	screenWidth := screenresolution.GetPrimary().Width
	if screenWidth == 0 {
		slog.Info("failed to get screen resolution width")
	}

	return &App{
		Window:         new(app.Window),
		Th:             initTheme(),
		IsFullScreen:   atomic.Bool{},
		DevelopMode:    atomic.Bool{},
		DevelopModeEnh: atomic.Bool{},
		MobileMode:     atomic.Bool{},
		Rendering:      atomic.Bool{},
		State:          storage.LoadState(),
		Client:         client.NewClient(httpClient),
		ScreenWidth:    screenWidth,
	}
}

func initTheme() *material.Theme {
	collection := fonts.ParseFont(fonts.SteticaFonts)
	th := material.NewTheme()
	th.Shaper = text.NewShaper(
		text.WithCollection(collection),
	)
	th.Palette.ContrastBg = color.NRGBA{R: 33, G: 92, B: 150, A: 255}
	return th
}

func (a *App) AutoSave(ctx context.Context, interval time.Duration, f func() storage.State) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state := f()
			storage.SaveState(state)
			return
		case <-ticker.C:
			if !a.DevelopMode.Load() {
				return
			}
			state := f()
			storage.SaveState(state)
		}
	}
}

func (a *App) AutoSaveCancelHelper() {
	if a.AutoSaveCancel != nil {
		a.AutoSaveCancel()
		a.AutoSaveCancel = nil
	}
}

func (a *App) Run(v ViewUI) error {
	a.Window.Option(app.Title("ROASTING BEANS"))

	var ops op.Ops

	v.SetState(a.State)
	go v.SetFlags()

	for {
		switch e := a.Window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			v.PointerEvents(gtx)

			v.Run(gtx, a.Th)

			e.Frame(gtx.Ops)
		}
	}
}
