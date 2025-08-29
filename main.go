package main

import (
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"log"
)

func main() {
	m, err := whisper.New("ggml-small.bin")
	if err != nil {
		log.Fatal(err)
	}
	_ = m
}
