package ffmpeg

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
)

var ErrFfmpegNotFound = errors.New("ffmpeg not found")

// Detect returns the path to the ffmpeg binary
func Detect() (string, error) {
	ffmpegExe := appFilename()

	// if ffmpeg is in the environment, use it
	if path := os.Getenv("FFMPEG_BIN"); path != "" {
		return path, nil
	}

	// if ffmpeg is in the path, use it
	if path, err := exec.LookPath(ffmpegExe); err == nil {
		return path, nil
	}

	// if ffmpeg is beside the app, use it
	if _, err := os.Stat(ffmpegExe); err == nil {
		return ffmpegExe, nil
	}

	return "", ErrFfmpegNotFound

}

// appFilename add extension to the filename if os needs it
func appFilename() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	} else {
		return "ffmpeg"
	}
}
