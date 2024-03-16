package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/TOomaAh/media_tools/internal/core/config"
)

type InternalFolderChooser struct {
	window *fyne.Window
	ch     *chan string
	config *config.Config
}

func NewInternalFolderChooser(window *fyne.Window, config *config.Config) (*InternalFolderChooser, *chan string) {
	ch := make(chan string)
	return &InternalFolderChooser{window: window, ch: &ch, config: config}, &ch
}

func (f *InternalFolderChooser) ShowOpen() string {
	dialog.ShowFolderOpen(func(uc fyne.ListableURI, err error) {
		if err != nil {
			return
		}

		if uc == nil {
			return
		}

		*f.ch <- uc.Path()

	}, *f.window)
	return ""
}

func (f *InternalFolderChooser) ShowSave() string {
	return ""
}

func (f *InternalFolderChooser) LoadGUI() error {
	return nil
}
