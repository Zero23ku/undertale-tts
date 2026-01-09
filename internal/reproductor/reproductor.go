package reproductor

import (
	"bytes"
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
