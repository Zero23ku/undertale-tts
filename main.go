package main

import (
	"context"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"undertale-tts/internal/charselector"
	"undertale-tts/internal/chzzk"
	"undertale-tts/internal/tiktok"
	"undertale-tts/internal/twitch"
	"undertale-tts/internal/web"
	"undertale-tts/internal/youtube"
)

var version = "v0.0.1"
var updateTime = false

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Undertale Text to Speech - " + version)

	mainWindow.SetOnClosed(func() {
		cancel()
		os.Exit(0)
	})

	go func() {
		web.StartWebServer()
	}()
	twitch.AppReference = &mainApp
	twitch.CTX = ctx
	twitch.InitTwitchConfigWindows(mainApp)

	youtube.CTX = ctx
	youtube.AppReference = &mainApp
	youtube.InitYoutubeWindow(mainApp)

	tiktok.CTX = ctx
	tiktok.AppReference = &mainApp
	tiktok.InitTiktokWindow(mainApp)

	chzzk.CTX = ctx
	chzzk.AppReference = &mainApp
	chzzk.InitChzzkWindow(mainApp)

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
			container.NewHBox(
				twitch.TwitchConnectButton, youtube.ConnectYTButton, tiktok.ConnectTiktokButton, chzzk.ConnectChzzkButton,
			),
		),
	)

	mainWindow.ShowAndRun()
}
