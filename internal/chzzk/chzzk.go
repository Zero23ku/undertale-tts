package chzzk

import (
	"context"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dhkimxx/GoChzzkChatCrawler/crawler"

	"undertale-tts/internal/charselector"
	"undertale-tts/internal/common"
	"undertale-tts/internal/reproductor"
	"undertale-tts/internal/whitelist"
)

var chzzkWindowIsOpen = false
var CHZZKWindow fyne.Window
var CHZZKErrorWindow fyne.Window
var CHZZKAlertWindow fyne.Window
var ConnectChzzkButton *widget.Button
var AppReference *fyne.App
var CTX context.Context
var chzzkLiveId = ""

func connectToChzzk(liveid string, ctx context.Context) {
	// Create a new crawler client with callback handler
	crawlerClient := crawler.NewCrawlerClient(liveid, 1, func(msg crawler.ChzzkChatMessage) {
		if common.IsCommandActive && strings.HasPrefix(msg.Content, common.TTSCommand) {
			if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(msg.Nickname) {
				reproductor.Reproduce(msg.Content, charselector.CurrentCharacter)
			} else if !whitelist.IsWhitelistActive {
				reproductor.Reproduce(msg.Content, charselector.CurrentCharacter)
			}
		} else if !common.IsCommandActive {
			if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(msg.Nickname) {
				reproductor.Reproduce(msg.Content, charselector.CurrentCharacter)
			} else if !whitelist.IsWhitelistActive {
				reproductor.Reproduce(msg.Content, charselector.CurrentCharacter)
			}
		}
		time.Sleep(500 * time.Millisecond)
	})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				CHZZKErrorWindow.Show()
				EnableButton()
				return
			}
		}()
		SetConnected()
		err := crawlerClient.Run()
		if err != nil {
			CHZZKErrorWindow.Show()
			return
		}
	}()

}

func InitChzzkWindow(app fyne.App) {
	CHZZKWindow = app.NewWindow("chzzk integration")
	initChzzkErrorWindow(app)
	CHZZKWindow.SetOnClosed(func() {
		chzzkWindowIsOpen = false
	})

	chzzkChannelNameInput := widget.NewEntry()
	chzzkChannelNameInput.SetPlaceHolder("Enter chzzk live id")
	chzzkChannelNameInput.Resize(fyne.NewSize(100, chzzkChannelNameInput.MinSize().Height))

	chzzkSubmit := widget.NewButton("Connect", func() {
		chzzkLiveId = chzzkChannelNameInput.Text
		if chzzkLiveId == "" {
			initChzzkAlertWindow(*AppReference, "You must enter a livestream id")
		} else {
			connectToChzzk(chzzkLiveId, CTX)
			CHZZKWindow.Close()
		}
	})

	form := widget.NewForm(
		widget.NewFormItem("chzzk live's id", chzzkChannelNameInput),
	)

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		chzzkSubmit,
	)

	CHZZKWindow.SetContent(
		container.NewVBox(form, centeredButton),
	)
	CHZZKWindow.Resize(fyne.NewSize(400, CHZZKWindow.Canvas().Size().Height))

	CHZZKWindow.SetOnClosed(func() {
		CHZZKWindow = nil
	})

	ConnectChzzkButton = widget.NewButton("Connect to chzzk", func() {
		if CHZZKWindow == nil {
			InitChzzkWindow(*AppReference)
			CHZZKWindow.Show()

		} else {
			CHZZKWindow.Show()
		}
	})
}

func initChzzkErrorWindow(app fyne.App) {
	CHZZKErrorWindow = app.NewWindow("Error!")
	CHZZKErrorWindow.SetContent(widget.NewLabel("An error ocurrier while connecting to chzzk, please try again!"))
}

func initChzzkAlertWindow(app fyne.App, msg string) {
	CHZZKErrorWindow = app.NewWindow("Error!")
	CHZZKErrorWindow.SetContent(widget.NewLabel(msg))
}

func SetConnected() {
	fyne.Do(func() {
		ConnectChzzkButton.SetText("Connected")
		ConnectChzzkButton.Disable()
	})
}

func EnableButton() {
	fyne.Do(func() {
		ConnectChzzkButton.SetText("Connect to chzzk")
		ConnectChzzkButton.Enable()
	})
}
