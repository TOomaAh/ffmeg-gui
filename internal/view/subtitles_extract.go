package view

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TOomaAh/media_tools/internal/core/config"
	"golang.org/x/net/context"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitlesExtractView struct {
	window fyne.Window
	config *config.Config
	files  *[]string
}

func NewSubtitlesExtractView(config *config.Config, files *[]string) *SubtitlesExtractView {
	return &SubtitlesExtractView{config: config, files: files}
}

func isTextSubtitles(codecName string) bool {
	return codecName == "subrip" || codecName == "ass" || codecName == "ssa"
}

func (sev *SubtitlesExtractView) LoadGUI(getFiles func() []string) fyne.CanvasObject {

	extractButton := widget.NewButton("Extract", func() {
		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancelFn()
		fmt.Println("Extracting subtitles of the following files:", getFiles())
		for _, file := range getFiles() {
			data, err := ffprobe.ProbeURL(ctx, file)
			if err != nil {
				fmt.Println("Error probing file", file, err)
				continue
			}

			for _, stream := range data.Streams {
				fmt.Println(stream.CodecType)
				if stream.CodecType == "subtitle" {

					if !isTextSubtitles(stream.CodecName) {
						fmt.Println(stream.CodecName)
						continue
					}

					var language string
					// language to lowercase
					if stream.TagList["language"] != "" {
						language = strings.ToLower(stream.TagList["language"].(string))
					} else {
						language = "unknown"
					}

					fileWithoutExtension := file[:len(file)-len(filepath.Ext(file))]

					outputFileName := fmt.Sprintf("%s.%s.", fileWithoutExtension, language)

					// if forced
					if stream.Disposition.Forced == 1 {
						outputFileName += "forced."
					}

					// if hearing impaired
					if stream.Disposition.HearingImpaired == 1 {
						outputFileName += "sdh."
					}

					outputFileName += "srt"

					// if file exists show popup
					if _, err := os.Stat(outputFileName); err == nil {
						dialog := widget.NewLabel("File already exists: " + outputFileName)
						popup := widget.NewModalPopUp(dialog, sev.window.Canvas())
						popup.Show()
						continue
					}

					cmd := exec.Command("ffmpeg.exe", "-i", file, "-map", fmt.Sprintf("0:%d", stream.Index), outputFileName)
					// get the output
					out, err := cmd.CombinedOutput()

					// print the output
					fmt.Println(string(out))

					if err != nil {
						fmt.Println("Error extracting subtitle", err)
					}

				}
			}

		}
	})

	return container.NewVBox(
		extractButton,
	)

}
