package charselector

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"undertale-tts/internal/constants"
	"undertale-tts/internal/reproductor"
)

var TestVoiceButton *widget.Button
var CharSelectionSelect *widget.Select
var CharSelectionLabel *canvas.Text
var CurrentCharacter string

func InitTestingControllers(app fyne.App) {
	CharSelectionLabel = canvas.NewText("Character Selection", theme.Color(theme.ColorNameForeground))
	CharSelectionSelect = widget.NewSelect(constants.CHAR_LIST, func(value string) {
		CurrentCharacter = value
	})

	TestVoiceButton = widget.NewButton("Test voice", func() {
		text := "Lorem ipsum dolor sit amet, consectetur"
		if CurrentCharacter == "" {
			//TODO
		} else {
			reproductor.Reproduce(text, CurrentCharacter)
		}
	})
}
