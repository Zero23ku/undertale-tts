package local

import (
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"undertale-tts/internal/constants"
	"undertale-tts/internal/reproductor"
)

var readButton *widget.Button
var textArea *widget.Entry
var saveButton *widget.Button
var filenameInput *widget.Entry
var folderButton *widget.Button
var saveDir string

var LocalContainer *fyne.Container
var LocalWindow fyne.Window
var LocalWindowIsOpen = false
var LocalButton *widget.Button
var AppReference *fyne.App
var LocalErrorWindow fyne.Window
var LocalSuccessWindow fyne.Window

var TestVoiceButton *widget.Button
var CharSelectionSelect *widget.Select
var CharSelectionLabel *canvas.Text
var CurrentCharacter string

const format = ".wav"

func InitLocalWindow(app fyne.App) {

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

	charSelectorContainer := container.NewCenter(
		container.NewHBox(
			container.NewVBox(
				CharSelectionLabel,
				CharSelectionSelect,
			),
			TestVoiceButton,
		),
	)

	LocalWindow = app.NewWindow("Read text offline")
	LocalWindow.SetOnClosed(func() {
		LocalWindowIsOpen = false
	})

	textArea = widget.NewMultiLineEntry()
	textArea.SetPlaceHolder("Enter text you want to be read")

	filenameInput = widget.NewEntry()
	filenameInput.SetPlaceHolder("Enter file name")

	folderButton = widget.NewButton("Select Folder", func() {
		dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				initErrorWindow(app, "Error selecting folder: "+err.Error())
				LocalErrorWindow.Show()
				return
			}
			if list == nil {
				folderButton.SetText("Select Folder")
				saveDir = ""
				return
			}
			saveDir = list.Path()
			folderButton.SetText("Folder: " + saveDir)
		}, LocalWindow).Show()
	})

	saveButton = widget.NewButton("Save audio", func() {
		if filenameInput.Text == "" {
			initErrorWindow(app, "You must enter a filename")
			LocalErrorWindow.Show()
			return
		}

		if saveDir == "" {
			initErrorWindow(app, "You must select a folder")
			LocalErrorWindow.Show()
			return
		}

		path := filepath.Join(saveDir, filenameInput.Text+format)
		reproductor.SaveAsWav(textArea.Text, CurrentCharacter, path)
		initSuccessWindow(app, "Audio clip saved!")
		LocalSuccessWindow.Show()
	})

	readButton = widget.NewButton("Read", func() {
		lines := strings.Split(textArea.Text, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			reproductor.Reproduce(line, CurrentCharacter)

		}

	})

	controls := container.NewVBox(
		charSelectorContainer,
		container.NewHBox(
			readButton,
			container.NewGridWrap(fyne.NewSize(300, filenameInput.MinSize().Height), filenameInput),
			folderButton,
			saveButton,
		),
	)

	LocalContainer = container.NewBorder(
		controls,
		nil, nil, nil,
		textArea,
	)
	LocalWindow.SetContent(LocalContainer)
	LocalWindow.Resize(fyne.NewSize(700, 500))

	LocalWindow.SetOnClosed(func() {
		LocalWindow = nil
	})

	LocalButton = widget.NewButton("Read text offline", func() {
		if LocalWindow == nil {
			InitLocalWindow(*AppReference)
			LocalWindow.Show()
		} else {
			LocalWindow.Show()
		}

	})

}

func initErrorWindow(app fyne.App, msg string) {
	LocalErrorWindow = app.NewWindow("Error!")
	LocalErrorWindow.SetContent(widget.NewLabel(msg))

}

func initSuccessWindow(app fyne.App, msg string) {
	LocalSuccessWindow = app.NewWindow("Success!")
	LocalSuccessWindow.SetContent(widget.NewLabel(msg))
}
