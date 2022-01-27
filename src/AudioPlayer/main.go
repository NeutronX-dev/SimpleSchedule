package AudioPlayer

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type AudioPlayer struct {
	AudioPath string
	Stream    beep.StreamSeekCloser
	Format    beep.Format
}

func (Player *AudioPlayer) Play() bool {
	finished := make(chan bool)
	Player.Stream.Seek(0)
	speaker.Play(beep.Seq(Player.Stream, beep.Callback(func() {
		finished <- true
	})))
	return <-finished
}

func (Player *AudioPlayer) Close() {
	Player.Stream.Close()
}

func New(path string) (*AudioPlayer, error) {
	f, err := os.Open(path)
	if err != nil {
		return &AudioPlayer{}, nil
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return &AudioPlayer{}, nil
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	return &AudioPlayer{AudioPath: path, Stream: streamer, Format: format}, nil
}
