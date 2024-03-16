package condition

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type AudioLanguage struct {
	ConditionChoice
}

func NewAudioLanguageCondition() *AudioLanguage {
	return &AudioLanguage{}
}

func (c *AudioLanguage) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is audio stream and language matches
		if s.CodecType == "audio" && c.checkString(s.Tags.Language) {
			return true
		}
	}
	return false
}

func (c *AudioLanguage) Name() string {
	return "Audio Language"
}

func (c *AudioLanguage) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *AudioLanguage) New() Condition {
	return &AudioLanguage{
		ConditionChoice{
			value: c.value,
		},
	}

}

func (c *AudioLanguage) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *AudioLanguage) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.value = s
	}
	return entry
}
