package condition

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type Container struct {
	ConditionChoice
}

func NewContainerCondition() *Container {
	return &Container{}
}

func (c *Container) Check(data *ffprobe.ProbeData) bool {
	// get extension of the file
	return filepath.Ext(data.Format.Filename)[1:] == c.value
}

func (c *Container) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *Container) Name() string {
	return "Container"
}

func (c *Container) New() Condition {
	return &Container{
		ConditionChoice{
			value: c.value,
		},
	}
}

func (c *Container) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *Container) GetEntry() fyne.Widget {
	selectWidget := widget.NewSelect([]string{
		"mp4",
		"mkv",
		"avi",
		"mov",
		"ts",
		"m4v",
	}, func(s string) {
		c.value = s
	})
	return selectWidget
}
