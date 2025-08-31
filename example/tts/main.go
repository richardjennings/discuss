package main

import (
	"github.com/gordonklaus/portaudio"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"log"
	"os"
	"runtime"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("expected text as a argument")
	}

	text := os.Args[1]

	sid := 0 // speaker id for multi speaker models

	config := sherpa.OfflineTtsConfig{}
	config.Model.Matcha.AcousticModel = "example/tts/model/matcha-icefall-en_US-ljspeech/model-steps-3.onnx"
	config.Model.Matcha.Vocoder = "example/tts/model/matcha-icefall-en_US-ljspeech/vocos-22khz-univ.onnx"
	config.Model.Matcha.Tokens = "example/tts/model/matcha-icefall-en_US-ljspeech/tokens.txt"
	config.Model.Matcha.DataDir = "example/tts/model/matcha-icefall-en_US-ljspeech/espeak-ng-data"
	config.Model.NumThreads = runtime.NumCPU()
	config.Model.Provider = "cpu"
	config.MaxNumSentences = 1
	config.Model.Debug = 1

	tts := sherpa.NewOfflineTts(&config)
	if tts == nil {
		log.Fatal("Could not create tts")
	}
	defer sherpa.DeleteOfflineTts(tts)

	audio := tts.Generate(text, sid, 1.0)

	_ = portaudio.Initialize()
	defer func() { _ = portaudio.Terminate() }()

	out := make([]float32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(audio.SampleRate), len(out), &out)
	e(err)
	defer func() { _ = stream.Close() }()
	e(stream.Start())
	defer func() { _ = stream.Stop() }()

	for i := 0; i < len(audio.Samples); i += 8192 {
		if i > len(audio.Samples) {
			return
		}
		ii := i + 8192
		if ii > len(audio.Samples) {
			ii = len(audio.Samples) - 1
		}
		out = audio.Samples[i:ii]
		_ = stream.Write()
	}

}

func e(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
