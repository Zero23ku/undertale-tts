package reproductor

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	"undertale-tts/internal/audios"
	"undertale-tts/internal/constants"
)

var (
	speakerInitOnce sync.Once
	speakerInitErr  error
)

func initSpeaker(format beep.Format) error {
	speakerInitOnce.Do(func() {
		speakerInitErr = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/14))
	})
	return speakerInitErr
}

func Reproduce(msg string, character string) {
	if character == constants.CHAR_SANS {
		ms := 69.99
		reader := bytes.NewReader(audios.Sans)
		streamer, format, err := wav.Decode(reader)

		if err != nil {
			//TODO
		}
		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		finalDuration := time.Duration(ms * float64(time.Millisecond))

		err = initSpeaker(format)
		if err != nil {
			//TODO
		}

		for _, char := range msg {
			if char == ' ' {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			sectionToReproduce := format.SampleRate.N(finalDuration)
			shot := buffer.Streamer(0, sectionToReproduce)
			resampled := beep.ResampleRatio(4, 1.0, shot)
			done := make(chan bool)

			speaker.Play(beep.Seq(resampled, beep.Callback(func() {
				done <- true
			})))

			<-done
		}
	}
}

func SaveAsWav(lines string, character string, fileName string) {
	if character == constants.CHAR_SANS {
		ms := 69.99
		reader := bytes.NewReader(audios.Sans)
		streamer, format, err := wav.Decode(reader)

		if err != nil {
			//TODO
		}
		buffer := beep.NewBuffer(format)
		buffer.Append(streamer)
		finalDuration := time.Duration(ms * float64(time.Millisecond))

		err = initSpeaker(format)
		if err != nil {
			//TODO
		}

		var parts []beep.Streamer
		fmt.Println(lines)
		sLines := strings.Split(lines, "\n")
		for _, line := range sLines {
			line = strings.TrimSpace(line)
			fmt.Println("linea", line)
			if line == "" {
				silence := beep.Silence(format.SampleRate.N(500 * time.Millisecond))
				parts = append(parts, silence)
			}

			for char, _ := range line {
				fmt.Println("letra", char)
				if char == ' ' {
					silence := beep.Silence(format.SampleRate.N(50 * time.Millisecond))
					parts = append(parts, silence)
				}
				sectionToReproduce := format.SampleRate.N(finalDuration)
				shot := buffer.Streamer(0, sectionToReproduce)
				resampled := beep.ResampleRatio(4, 1.0, shot)
				parts = append(parts, resampled)

			}

		}

		finalStreamer := beep.Seq(parts...)
		f, err := os.Create(fileName)
		if err != nil {
			return
		}
		defer f.Close()
		err = wav.Encode(f, finalStreamer, format)
		if err != nil {
			return
		}
	}

}
