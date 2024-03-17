package view

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TOomaAh/media_tools/internal/component"
	"github.com/TOomaAh/media_tools/internal/core"
)

type FiltersView struct {
	container  *fyne.Container
	choices    *[]*component.ConditionalWidget
	workerPool *core.Dispatcher
}

func NewFiltersView() *FiltersView {
	return &FiltersView{
		container: container.NewVBox(),
		choices:   &[]*component.ConditionalWidget{},
	}
}

func (fv *FiltersView) LoadGUI(files *[]string, progressChan *chan float64) fyne.CanvasObject {

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

	filtersButton := widget.NewButton("Filters", fv.Filter(progressChan, files))

	return container.NewBorder(filtersButton,
		container.NewHBox(minusButton, plusButton), nil, nil, fv.container)

}

func (fv *FiltersView) Filter(progressChan *chan float64, files *[]string) func() {

	return func() {
		fv.workerPool = core.NewDispatcher(5)

		fv.workerPool.Run()
		output := []string{}

		for i, file := range *files {
			*progressChan <- float64(i) / float64(len(*files))

			if fileInfo, err := os.Stat(file); err == nil && fileInfo.IsDir() {
				processFolder(fv.workerPool, file, &output, fv.choices)
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
							output = append(output, file)
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
func processFolder(workPool *core.Dispatcher, folderPath string, output *[]string, choices *[]*component.ConditionalWidget) {
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
					*output = append(*output, path)
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
