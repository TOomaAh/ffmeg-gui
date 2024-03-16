package view

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/TOomaAh/media_tools/internal/core/config"
)

type ConfigView struct {
	app    *fyne.App
	config *config.Config
}

func NewConfigView(app *fyne.App, config *config.Config) *ConfigView {
	return &ConfigView{app: app, config: config}
}

func (c *ConfigView) LoadGUI() error {
	w := (*c.app).NewWindow("Config")
	w.Resize(fyne.NewSize(400, 400))

	extensionInput := widget.NewEntry()
	extensionInput.SetPlaceHolder("Extension")

	extensionInput.Text = strings.Join(c.config.AcceptedExtensions, ", ")

	saveButton := widget.NewButton("Save", func() {
		c.config.AcceptedExtensions = strings.Split(extensionInput.Text, ",")
		c.config.Update()
		w.Close()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Accepted extensions"),
		extensionInput,
		saveButton,
	))

	w.Show()
	return nil
}
