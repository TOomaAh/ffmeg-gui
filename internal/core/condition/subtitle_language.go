package condition

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitleLanguage struct {
	ConditionChoice
}

func NewSubtitleLanguageCondition() *SubtitleLanguage {
	return &SubtitleLanguage{}
}

func (c *SubtitleLanguage) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is subtitle stream and language matches
		if s.CodecType == "subtitle" && c.checkString(s.Tags.Language) {
			return true
		}
	}
	return false
}

func (c *SubtitleLanguage) Name() string {
	return "Subtitle language"
}

func (c *SubtitleLanguage) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *SubtitleLanguage) New() Condition {
	return &SubtitleLanguage{
		ConditionChoice{
			value: c.value,
		},
	}
}

func (c *SubtitleLanguage) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *SubtitleLanguage) SetValue(value string) {
	c.value = value
}

func (c *SubtitleLanguage) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.value = s
	}
	return entry
}
