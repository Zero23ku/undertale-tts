package main

import (
	"context"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"undertale-tts/internal/charselector"
	"undertale-tts/internal/chzzk"
	"undertale-tts/internal/common"
	"undertale-tts/internal/local"
	"undertale-tts/internal/stt"
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

	stt.AppReference = &mainApp
	stt.InitSTTWindow(mainApp)

	common.InitCommandCheck()
	common.InitKofiButton()

	local.AppReference = &mainApp
	local.InitLocalWindow(mainApp)

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

	commandContent := container.NewCenter(
		container.NewVBox(
			container.NewHBox(
				common.ActivateCommand,
				common.InputCommand,
			),
		),
	)

	footer := container.NewVBox(common.KofiButton)

	mainContent := container.NewVBox(
		charSelectorContainer,
		container.NewVBox(
			container.NewHBox(
				twitch.TwitchConnectButton, youtube.ConnectYTButton, tiktok.ConnectTiktokButton, chzzk.ConnectChzzkButton,
			),
			commandContent,
		),
	)

	mainWindow.SetContent(
		container.New(
			layout.NewBorderLayout(nil, footer, nil, nil),
			footer,
			container.NewVBox(mainContent, local.LocalButton, stt.STTButton),
		),
	)

	mainWindow.ShowAndRun()
}
