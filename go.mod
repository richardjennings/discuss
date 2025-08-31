module github.com/richardjennings/discuss

go 1.24.5

replace github.com/ggerganov/whisper.cpp/bindings/go => ./whisper.cpp/bindings/go

require (
	github.com/ggerganov/whisper.cpp/bindings/go v0.0.0-00010101000000-000000000000
	github.com/go-audio/wav v1.1.0
	github.com/gordonklaus/portaudio v0.0.0-20250206071425-98a94950218b
)

require (
	github.com/go-audio/audio v1.0.0 // indirect
	github.com/go-audio/riff v1.0.0 // indirect
)
