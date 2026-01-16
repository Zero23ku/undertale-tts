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

var UpdateButton *widget.Button
var DocsButton *widget.Button

var IsCommandActive = false
var TTSCommand = "!tts"
var kofiUrl *url.URL
var githubUrl *url.URL
var docsUrl *url.URL

func InitCommandCheck() {
	ActivateCommand = widget.NewCheck("Use Command", func(value bool) {
		IsCommandActive = value
	})
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

func InitUpdateButton() {

	githubUrl = &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/Zero23ku/undertale-tts/releases",
	}

	UpdateButton = widget.NewButton("New Version Avaible", func() {
		fyne.CurrentApp().OpenURL(githubUrl)
	})

}

func InitDocsButton() {
	docsUrl = &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/Zero23ku/undertale-tts/blob/main/docs/docs.md",
	}

	DocsButton = widget.NewButton("How to use", func() {
		fyne.CurrentApp().OpenURL(docsUrl)
	})
}
