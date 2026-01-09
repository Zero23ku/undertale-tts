package main

import (
	"context"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"undertale-tts/internal/charselector"
)

var version = "v0.0.1"
var updateTime = false

func main() {
	_, cancel := context.WithCancel(context.Background())
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Undertale Text to Speech - " + version)

	mainWindow.SetOnClosed(func() {
		cancel()
		os.Exit(0)
	})

	charselector.InitTestingControllers(mainApp)

	charSelectorContainer := container.NewCenter(
		container.NewHBox(
			container.NewVBox(
				charselector.CharSelectionLabel,
				charselector.CharSelectionSelect,
			),
			charselector.TestVoiceButton,
		),
	)

	mainWindow.SetContent(
		container.NewVBox(

			charSelectorContainer,
		),
	)

	mainWindow.ShowAndRun()
}
