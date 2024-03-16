package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/TOomaAh/media_tools/internal/core/config"
)

type FileChooser interface {
	ShowOpen() string
	ShowSave() string
}

type InternalFileChooser struct {
	window *fyne.Window
	ch     *chan string
	dialog *dialog.FileDialog
}

func NewInternalFileChooser(window *fyne.Window, config *config.Config) (*InternalFileChooser, *chan string) {
	ch := make(chan string)
	dialog := dialog.NewFileOpen(
		func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				return
			}

			if uc == nil {
				return
			}

			ch <- uc.URI().Path()
		}, *window)
	return &InternalFileChooser{ch: &ch, dialog: dialog, window: window}, &ch
}

func (f *InternalFileChooser) ShowOpen() {
	f.dialog.Show()
}
