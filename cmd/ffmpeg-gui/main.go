package main

import (
	"github.com/TOomaAh/media_tools/internal/core/config"
	"github.com/TOomaAh/media_tools/internal/ffmpeg"
	"github.com/TOomaAh/media_tools/internal/media_tools"
)

func main() {
	var Config = config.NewConfig()
	ffmpegPath, err := ffmpeg.Detect()

	if err != nil {
		panic(err)
	}

	media_tools.Run(ffmpegPath, Config)

}
