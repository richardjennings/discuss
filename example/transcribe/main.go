package main

import (
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"io"
	"time"

	"github.com/gordonklaus/portaudio"
)

func main() {

	// load whisper.cpp model
	model, err := whisper.New("ggml-small.bin")
	e(err)

	fmt.Println("press ctrl-c to exit")
	for {
		// try to use microphone to read raw pcm
		data, err := Listen()
		e(err)

		// use whisper to generate text from pcm audio
		text, err := STT(model, data)
		e(err)

		fmt.Println(string(text))
	}
}

func Listen() ([]float32, error) {
	var data []float32
	done := make(chan bool)

	// Read audio device into data until <-done
	go func() {
		_ = portaudio.Initialize()
		time.Sleep(1)
		defer func() { _ = portaudio.Terminate() }()
		in := make([]float32, 64)
		stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
		e(err)
		defer func() { _ = stream.Close() }()
		e(stream.Start())
		for {
			select {
			case <-done:
				break
			default:
				e(stream.Read())
				data = append(data, in...)
			}
		}
	}()

	fmt.Println("Speak, press enter to stop.")
	var input []byte
	_, _ = fmt.Scanln(&input)
	done <- true
	return data, nil
}

func STT(model whisper.Model, data []float32) ([]byte, error) {
	context, err := model.NewContext()
	if err != nil {
		return nil, err
	}
	var cb whisper.SegmentCallback
	if err := context.Process(data, nil, cb, func(i int) {}); err != nil {
		return nil, err
	}
	return text(context)
}

func text(context whisper.Context) ([]byte, error) {
	var text []byte
	for {
		segment, err := context.NextSegment()
		if err == io.EOF {
			return text, nil
		} else if err != nil {
			return nil, err
		}
		text = append(text, []byte(segment.Text)...)
	}
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
