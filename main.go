package main

import (
	"log"
	"os"
	a "voicer/internal/app"
	"voicer/internal/logger"
	"voicer/internal/ui"

	"gioui.org/app"
)

func main() {
	logger.NewLogger()

	go func() {
		application := a.NewApp()

		u := ui.NewUI(application)

		if err := application.Run(u); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
