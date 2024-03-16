package condition

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type Bitrate struct {
	ConditionChoice
}

func NewBitrateCondition() *Bitrate {
	return &Bitrate{}
}

func (c *Bitrate) Check(data *ffprobe.ProbeData) bool {
	return c.checkInt(data.Format.BitRate)
}

func (c *Bitrate) Name() string {
	return "Bitrate"
}

func (c *Bitrate) GetPossibleConditions() []string {
	return []string{"equals", "greater than", "less than", "greater or equals", "less or equals"}
}

func (c *Bitrate) New() Condition {
	return &Bitrate{
		ConditionChoice{
			value: c.value,
		},
	}
}

func (c *Bitrate) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *Bitrate) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.value = s
	}
	return entry
}
