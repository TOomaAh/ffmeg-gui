package media_tools

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TOomaAh/media_tools/internal/component"
	"github.com/TOomaAh/media_tools/internal/core/config"
	"github.com/TOomaAh/media_tools/internal/view"
	"github.com/kbinani/screenshot"
)

var (
	selectedFiles     = []string{}
	selectedFilesText = "%d files selected"
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
			container.NewTabItem("Extract", view.NewSubtitlesExtractView(config, &selectedFiles).LoadGUI(func() []string {
				return selectedFiles
			})),
		)),
		container.NewTabItem("Filters", view.NewFiltersView().LoadGUI(&selectedFiles, &ProgressBarChannel)),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	folder, chanFolder := component.NewInternalFolderChooser(&w, config)

	file, chanFile := component.NewInternalFileChooser(&w, config)

	// show number of selected files
	selectedFileText := widget.NewLabel(fmt.Sprintf(selectedFilesText, len(selectedFiles)))

	progressBar := widget.NewProgressBar()

	go func() {
		for {
			select {
			case folder := <-*chanFolder:
				selectedFiles = append(selectedFiles, folder)
				selectedFileText.SetText(fmt.Sprintf(selectedFilesText, len(selectedFiles)))
			case file := <-*chanFile:
				selectedFiles = append(selectedFiles, file)
				selectedFileText.SetText(fmt.Sprintf(selectedFilesText, len(selectedFiles)))
			case progress := <-ProgressBarChannel:
				fmt.Println("Progress", progress)
				progressBar.SetValue(float64(progress))
			}
		}
	}()

	fileButton := widget.NewButton("Choose file", func() {
		file.ShowOpen()
	})

	folderButton := widget.NewButton("Choose folder", func() {
		folder.ShowOpen()
	})

	grid := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewVBox(
			container.NewHBox(
				selectedFileText,
				fileButton,
				folderButton,
				progressBar,
			),
			tabs,
		),
	)

	w.SetContent(grid)

	w.ShowAndRun()

}
