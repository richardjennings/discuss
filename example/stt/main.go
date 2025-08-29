package main

import (
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	model, err := whisper.New("ggml-small.bin")
	if err != nil {
		log.Fatal(err)
	}

	// Create processing context
	context, err := model.NewContext()
	e(err)
	fmt.Printf("\n%s\n", context.SystemInfo())

	var data []float32

	path := "whisper.cpp/samples/jfk.wav"
	fmt.Printf("loading %s\n", path)
	fh, err := os.Open(path)
	e(err)
	defer func() { _ = fh.Close() }()

	// Decode the WAV file - load the full buffer
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		e(err)
	} else if dec.SampleRate != whisper.SampleRate {
		e(fmt.Errorf("unsupported sample rate: %d", dec.SampleRate))
	} else if dec.NumChans != 1 {
		e(fmt.Errorf("unsupported number of channels: %d", dec.NumChans))
	} else {
		data = buf.AsFloat32Buffer().Data
	}

	context.ResetTimings()

	var cb whisper.SegmentCallback
	err = context.Process(data, nil, cb, nil)
	e(err)

	context.PrintTimings()

	e(Output(os.Stdout, context))
}

func e(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Output(w io.Writer, context whisper.Context) error {
	for {
		segment, err := context.NextSegment()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		fmt.Fprintf(w, "[%6s->%6s]", segment.Start.Truncate(time.Millisecond), segment.End.Truncate(time.Millisecond))
		fmt.Fprintln(w, " ", segment.Text)
	}
}
