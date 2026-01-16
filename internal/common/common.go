package common

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"undertale-tts/internal/assets"
)

var KofiButton *widget.Button
var ActivateCommand *widget.Check
var InputCommand *widget.Entry

var IsCommandActive = false
var Pitch = 1.2
var TTSCommand = "!tts"
var kofiUrl *url.URL

func InitCommandCheck() {
	ActivateCommand = widget.NewCheck("Use Command", func(value bool) {
		IsCommandActive = value
	})
}

func InitCommandInput() {
	InputCommand = widget.NewEntry()
	InputCommand.Text = TTSCommand
}

func InitKofiButton() {
	kofiUrl = &url.URL{
		Scheme: "https",
		Host:   "ko-fi.com",
		Path:   "/I3I41O6OUD",
	}

	res := fyne.NewStaticResource("cup-border.png", assets.Cup)

	KofiButton = widget.NewButtonWithIcon("Support me!", res, func() {
		fyne.CurrentApp().OpenURL(kofiUrl)
	})
}
