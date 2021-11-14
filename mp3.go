package mp3

import (
	"golang.org/x/net/context"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

func Play(ctx context.Context, name string) {
	f, err := os.Open(name)
	if err != nil {
		return
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return
	}
	defer streamer.Close()

	done := make(chan bool)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			speaker.Close()
			return
		}
	}
}
