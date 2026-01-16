package stt

import (
	"bufio"
	"context"
	"os"
	"regexp"
	"strconv"
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

var initSttButton *widget.Button
var fileButton *widget.Button

var STTContainer *fyne.Container
var STTWindow fyne.Window
var STTWindowIsOpen = false
var STTButton *widget.Button
var AppReference *fyne.App
var STTErrorWindow fyne.Window
var file *os.File
var cancel context.CancelFunc
var fileName string

var TestVoiceButton *widget.Button
var CharSelectionSelect *widget.Select
var CharSelectionLabel *canvas.Text
var CurrentCharacter string

func InitSTTWindow(app fyne.App) {

	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())

	STTWindow = app.NewWindow("Read from STT")
	STTWindow.SetOnClosed(func() {
		STTWindowIsOpen = false
		cancel()
	})

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

	fileButton = widget.NewButton("Select text file", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				initErrorWindow(app, err.Error())
				STTErrorWindow.Show()
				return
			}

			if reader == nil {
				file = nil
				fileButton.SetText("Select text file")
				return
			}

			defer reader.Close()

			file, err = os.Open(reader.URI().Path())
			if err != nil {
				STTErrorWindow.Show()
				initErrorWindow(app, "Error trying to open file: "+err.Error())
			}

			fileName = file.Name()
			fileButton.SetText(fileName)

		}, STTWindow).Show()
	})

	initSttButton = widget.NewButton("Read STT file", func() {
		if file == nil {
			initErrorWindow(app, "Select a file first")
			STTErrorWindow.Show()
			return
		}
		initSttButton.SetText("Reading...")
		initSttButton.Disable()
		go func(ctx context.Context) {
			reader := bufio.NewReader(file)
			timeRegex := regexp.MustCompile(`\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}`)
			for {

				select {
				case <-ctx.Done():
					return
				default:
					line, err := reader.ReadString('\n')
					if err != nil {
						time.Sleep(500 * time.Millisecond)
						continue
					}

					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}

					if _, err := strconv.Atoi(line); err == nil {
						continue
					}

					if timeRegex.MatchString(line) {
						continue
					}
					reproductor.Reproduce(line, CurrentCharacter)

					time.Sleep(1500 * time.Millisecond)
				}

			}
		}(ctx)
	})

	controls := container.NewVBox(
		charSelectorContainer,
		fileButton,
		initSttButton,
	)

	STTWindow.SetContent(controls)
	STTWindow.Resize(fyne.NewSize(700, 300))

	STTWindow.SetOnClosed(func() {
		STTWindow = nil
	})

	STTButton = widget.NewButton("Read text from STT", func() {

		if STTWindow == nil {
			InitSTTWindow(*AppReference)
			STTWindow.Show()
		} else {
			STTWindow.Show()
		}

	})
}

func initErrorWindow(app fyne.App, msg string) {
	STTErrorWindow = app.NewWindow("Error!")
	STTErrorWindow.SetContent(widget.NewLabel(msg))

}
