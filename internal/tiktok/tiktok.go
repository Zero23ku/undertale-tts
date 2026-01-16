package tiktok

import (
	"context"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/steampoweredtaco/gotiktoklive"

	"undertale-tts/internal/charselector"
	"undertale-tts/internal/common"
	"undertale-tts/internal/reproductor"
	"undertale-tts/internal/whitelist"
)

var TiktokWindow fyne.Window
var TiktokErrorWindow fyne.Window
var ConnectTiktokButton *widget.Button
var AppReference *fyne.App
var CTX context.Context
var tiktokChannelName = ""

func connectToTikTokChat(username string, ctx context.Context) {
	tiktok, err := gotiktoklive.NewTikTok()
	if err != nil {
		//logging.CreateLog("Error initializing Tiktok..", err)
		log.Fatal(err)
		return
	}
	live, err := tiktok.TrackUser(username)

	if err != nil {
		TiktokErrorWindow.Show()
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				EnableButton()
				return
			case event := <-live.Events:
				switch e := event.(type) {
				case gotiktoklive.ChatEvent:
					tiktokMsg := gotiktoklive.ChatEvent(e).Comment
					chatter := gotiktoklive.ChatEvent(e).User.Nickname
					if common.IsCommandActive && strings.HasPrefix(tiktokMsg, common.TTSCommand) {
						if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(chatter) {
							reproductor.Reproduce(tiktokMsg, charselector.CurrentCharacter)
						} else if !whitelist.IsWhitelistActive {
							reproductor.Reproduce(tiktokMsg, charselector.CurrentCharacter)
						}

					} else if !common.IsCommandActive {
						if whitelist.IsWhitelistActive && whitelist.IsUserInWhitelist(chatter) {
							reproductor.Reproduce(tiktokMsg, charselector.CurrentCharacter)
						} else if !whitelist.IsWhitelistActive {
							reproductor.Reproduce(tiktokMsg, charselector.CurrentCharacter)
						}
					}

					time.Sleep(time.Duration(1200) * time.Millisecond)
				}
			}
		}
	}()

	SetConnected()
}

func InitTiktokWindow(app fyne.App) {
	TiktokWindow = app.NewWindow("Tiktok integration")
	initTiktokErrorWindow(app)

	tiktokChannelNameInput := widget.NewEntry()
	tiktokChannelNameInput.SetPlaceHolder("Enter Tiktok channel's name")
	tiktokChannelNameInput.Resize(fyne.NewSize(100, tiktokChannelNameInput.MinSize().Height))

	tiktokSubmit := widget.NewButton("Connect", func() {
		tiktokChannelName = tiktokChannelNameInput.Text
		connectToTikTokChat(tiktokChannelName, CTX)
		TiktokWindow.Close()
	})

	form := widget.NewForm(
		widget.NewFormItem("Tiktok Channel's name", tiktokChannelNameInput),
	)

	centeredButton := container.New(
		layout.NewBorderLayout(nil, nil, layout.NewSpacer(), layout.NewSpacer()),
		tiktokSubmit,
	)

	TiktokWindow.SetContent(container.NewVBox(form, centeredButton))
	TiktokWindow.Resize(fyne.NewSize(400, 100))

	TiktokWindow.SetOnClosed(func() {
		TiktokWindow = nil
	})

	ConnectTiktokButton = widget.NewButton("Connect to Tiktok", func() {
		if TiktokWindow == nil {
			InitTiktokWindow(*AppReference)
			TiktokWindow.Show()
		} else {
			TiktokWindow.Show()
		}

	})

}

func initTiktokErrorWindow(app fyne.App) {
	TiktokErrorWindow = app.NewWindow("Error!")
	TiktokErrorWindow.SetContent(widget.NewLabel("An error ocurried while connecting to Tiktok, please try again."))

}

func SetConnected() {
	fyne.Do(func() {
		ConnectTiktokButton.SetText("Connected")
		ConnectTiktokButton.Disable()
	})
}

func EnableButton() {
	fyne.Do(func() {
		ConnectTiktokButton.SetText("Connect to Tiktok")
		ConnectTiktokButton.Enable()
	})
}
