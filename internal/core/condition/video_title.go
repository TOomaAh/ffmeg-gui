package condition

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type VideoTitle struct {
	ConditionChoice
}

func NewVideoTitleCondition() *VideoTitle {
	return &VideoTitle{}
}

func (c *VideoTitle) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is video stream and title matches
		if s.CodecType == "video" && c.checkString(s.Tags.Title) {
			return true
		}
	}
	return false
}

func (c *VideoTitle) Name() string {
	return "Video title"
}

func (c *VideoTitle) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *VideoTitle) New() Condition {
	return &VideoTitle{
		ConditionChoice{
			value: c.value,
		},
	}
}

func (c *VideoTitle) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *VideoTitle) SetValue(value string) {
	c.value = value
}

func (c *VideoTitle) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.value = s
	}
	return entry
}
