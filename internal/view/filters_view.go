package view

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TOomaAh/media_tools/internal/component"
	"github.com/TOomaAh/media_tools/internal/core"
	"github.com/TOomaAh/media_tools/internal/core/config"
)

var (
	selectedFilesText = "%d files selected"
)

type FiltersView struct {
	container  *fyne.Container
	choices    *[]*component.ConditionalWidget
	workerPool *core.Dispatcher
	w          *fyne.Window
	config     *config.Config
	files      *[]string
}

func NewFiltersView(w *fyne.Window, config *config.Config) *FiltersView {
	return &FiltersView{
		container: container.NewVBox(),
		choices:   &[]*component.ConditionalWidget{},
		w:         w,
		config:    config,
		files:     &[]string{},
	}
}

func (fv *FiltersView) LoadGUI() fyne.CanvasObject {

	folder, chanFolder := component.NewInternalFolderChooser(fv.w, fv.config)

	file, chanFile := component.NewInternalFileChooser(fv.w, fv.config)
	selectedFileText := widget.NewLabel(fmt.Sprintf(selectedFilesText, len(*fv.files)))

	progressBar := widget.NewProgressBar()

	go func() {
		for {
			select {
			case folder := <-*chanFolder:
				*fv.files = append(*fv.files, folder)
				selectedFileText.SetText(fmt.Sprintf(selectedFilesText, len(*fv.files)))
			case file := <-*chanFile:
				*fv.files = append(*fv.files, file)
				selectedFileText.SetText(fmt.Sprintf(selectedFilesText, len(*fv.files)))
			}
		}
	}()

	fileButton := widget.NewButton("Choose file", func() {
		file.ShowOpen()
	})

	folderButton := widget.NewButton("Choose folder", func() {
		folder.ShowOpen()
	})

	clearButton := widget.NewButton("Clear", func() {
		*fv.files = []string{}
		selectedFileText.SetText(fmt.Sprintf(selectedFilesText, len(*fv.files)))
	})

	plusButton := widget.NewButton("+", func() {
		newFilter := component.NewConditionalWidget(fv.container)
		fv.container.Add(newFilter)
		*fv.choices = append(*fv.choices, newFilter)
	})

	minusButton := widget.NewButton("-", func() {
		if len(*fv.choices) > 0 {
			fv.container.Remove((*fv.choices)[len(*fv.choices)-1])
			*fv.choices = (*fv.choices)[:len(*fv.choices)-1]
		}
	})

	conditionalWidget := component.NewConditionalWidget(fv.container)
	*fv.choices = append(*fv.choices, conditionalWidget)

	// add first conditional widget
	fv.container.Add(conditionalWidget)

	filtersButton := widget.NewButton("Filters", fv.Filter())

	return container.NewBorder(container.NewVBox(
		container.NewHBox(
			selectedFileText,
			fileButton,
			folderButton,
			clearButton,
			progressBar,
		),
		filtersButton,
	),
		container.NewHBox(minusButton, plusButton), nil, nil, fv.container)

}

func (fv *FiltersView) Filter() func() {

	return func() {
		fv.workerPool = core.NewDispatcher(5)

		fv.workerPool.Run()
		output := []string{}

		// init mutex
		var mutex sync.Mutex

		for _, file := range *fv.files {

			if fileInfo, err := os.Stat(file); err == nil && fileInfo.IsDir() {
				processFolder(fv.workerPool, file, &output, fv.choices, &mutex)
			} else {
				if err != nil {
					fmt.Println("Error getting file info", file, err)
					continue
				}

				job := core.Job{
					Work: func() (bool, error) {
						// do the work
						filer := core.NewFilter(file, fv.choices)
						if filer.Apply() {
							mutex.Lock()
							output = append(output, file)
							mutex.Unlock()
							return true, nil
						}
						return false, errors.New("file does not match the filters")
					},
				}
				fv.workerPool.AddJob(job)
			}

		}

		fv.workerPool.Close()

		// sort the output
		sort.Strings(output)

		// show the results

		resultWindow := fyne.CurrentApp().NewWindow("Result")
		resultWindow.Resize(fyne.NewSize(400, 400))
		if len(output) == 0 {
			resultWindow.SetContent(widget.NewLabel("No files found"))
			resultWindow.Show()
			return
		}

		resultList := widget.NewList(func() int {
			return len(output)
		}, func() fyne.CanvasObject {
			if len(output) > 0 {
				return widget.NewLabel(output[0])
			} else {
				return widget.NewLabel("No results")
			}
		}, func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(output[i])
		},
		)

		resultWindow.SetContent(resultList)

		resultWindow.Show()
	}

}

// processFolder processes the folder and all subfolder (recursive) and returns the files that match the filters
func processFolder(workPool *core.Dispatcher, folderPath string, output *[]string, choices *[]*component.ConditionalWidget, mutex *sync.Mutex) {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		job := core.Job{
			Work: func() (bool, error) {
				filer := core.NewFilter(path, choices)
				if filer.Apply() {
					mutex.Lock()

					*output = append(*output, path)
					mutex.Unlock()
					return true, nil
				}
				return false, errors.New("file does not match the filters")
			},
		}

		workPool.AddJob(job)

		return nil
	})

	if err != nil {
		fmt.Println("Error walking folder", folderPath, err)
	}
}
