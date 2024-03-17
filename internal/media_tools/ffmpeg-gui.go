package media_tools

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TOomaAh/media_tools/internal/core/config"
	"github.com/TOomaAh/media_tools/internal/view"
	"github.com/kbinani/screenshot"
)

var ProgressBarChannel = make(chan float64)

func Run(ffmpegPath string, config *config.Config) {

	a := app.New()
	w := a.NewWindow("Hello World")

	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Config", func() {
				view.NewConfigView(&a, config).LoadGUI()
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
				os.Exit(0)
			}),
		),
	)

	w.SetMainMenu(menu)

	// get screen size
	screenSize := screenshot.GetDisplayBounds(0)
	w.Resize(fyne.NewSize(float32(screenSize.Dx()/2), float32(screenSize.Dy()/2)))

	tabs := container.NewAppTabs(
		container.NewTabItem("Video", container.NewAppTabs(
			container.NewTabItem("Convert", widget.NewLabel("Convert")),
		)),
		container.NewTabItem("Audio", container.NewAppTabs(
			container.NewTabItem("Change Language", widget.NewLabel("Change language")),
		)),
		container.NewTabItem("Subtitles", container.NewAppTabs(
			container.NewTabItem("Convert", widget.NewLabel("Convert")),
			container.NewTabItem("Extract", view.NewSubtitlesExtractView(config, &[]string{}).LoadGUI(func() []string {
				return []string{}
			})),
		)),
		container.NewTabItem("Filters", view.NewFiltersView(&w, config).LoadGUI()),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	// show number of selected files

	grid := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewVBox(

			tabs,
		),
	)

	w.SetContent(grid)

	w.ShowAndRun()

}
